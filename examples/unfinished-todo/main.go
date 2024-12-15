package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/aquiffoo/goblet"
	_ "github.com/mattn/go-sqlite3"
)

const dbpath string = "./db/todo.sqlite3"

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Todos   []Todo      `json:"todos,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func getDB(dbpath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL
        );
        CREATE TABLE IF NOT EXISTS todos (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            description TEXT,
            done INTEGER DEFAULT 0,
            user_id INTEGER,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
    `)
	if err != nil {
		panic(err)
	}

	return db
}

var db *sql.DB = getDB(dbpath)

func login(username string, password string) (bool, string) {
	var storedHash string
	var userID int
	err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&userID, &storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, ""
		}
		fmt.Println("Error during login query:", err)
		return false, ""
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return false, ""
	}

	return true, strconv.Itoa(userID)
}

func register(username string, password string) (bool, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("error hashing password: %w", err)
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return false, fmt.Errorf("username already exists")
		}
		return false, fmt.Errorf("error registering user: %w", err)
	}

	return true, nil
}

func addTodo(title string, description string, userID int) (bool, error) {
	_, err := db.Exec("INSERT INTO todos (title, description, user_id) VALUES (?, ?, ?)", title, description, userID)
	if err != nil {
		return false, fmt.Errorf("error adding todo: %w", err)
	}

	return true, nil
}

func getTodos(userID int) ([]Todo, error) {
	rows, err := db.Query("SELECT id, title, description, done FROM todos WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error getting todos: %w", err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Done)
		if err != nil {
			return nil, fmt.Errorf("error scanning todo: %w", err)
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over todos: %w", err)
	}

	return todos, nil
}

func main() {
	app := goblet.New(true)

	app.Handle("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"title":   "Goblet Todo - Home",
			"lang":    "en",
			"charset": "utf-8",
			"success": false,
		}

		userID, loggedIn := getLoggedInUserID(r)
		if loggedIn {
			todos, err := getTodos(userID)
			if err != nil {
				data["error"] = "Failed to load todos"
			} else {
				data["success"] = true
				data["username"] = getUserName(userID)
				data["todos"] = todos
			}
		}

		app.Render(w, "index.html", data)
	})

	app.Handle("/login", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"title":   "Goblet Todo - Login",
			"lang":    "en",
			"charset": "utf-8",
		}
		app.Render(w, "login.html", data)
	})

	app.Handle("/register", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"title":   "Goblet Todo - Register",
			"lang":    "en",
			"charset": "utf-8",
		}
		app.Render(w, "register.html", data)
	})

	app.Handle("/api/login", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		success, userID := login(username, password)

		var response APIResponse
		if success {
			setSessionCookie(w, userID)

			response = APIResponse{
				Success: true,
				Message: "Login successful",
				Data: map[string]string{
					"user_id": userID,
				},
			}
		} else {
			response = APIResponse{
				Success: false,
				Message: "Invalid username or password",
			}
			w.WriteHeader(http.StatusUnauthorized)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	app.Handle("/api/register", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		success, err := register(username, password)

		var response APIResponse
		if success {
			response = APIResponse{
				Success: true,
				Message: "Registration successful",
			}
		} else {
			response = APIResponse{
				Success: false,
				Message: err.Error(),
			}
			w.WriteHeader(http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	app.Handle("/api/todos", func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		todos, err := getTodos(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := APIResponse{
			Success: true,
			Todos:   todos,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	app.Handle("/api/add_todo", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		title := r.FormValue("title")
		description := r.FormValue("description")
		userIDStr := r.FormValue("user_id")

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		success, err := addTodo(title, description, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := APIResponse{
			Success: success,
			Message: "Todo added successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	err := app.Serve("6969")
	if err != nil {
		panic(err)
	}
}

func setSessionCookie(w http.ResponseWriter, userID string) {
	sessionValue := userID

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionValue,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})
}

func getLoggedInUserID(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, false
	}

	userIDStr := cookie.Value
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, false
	}

	return userID, true
}

func getUserName(userID int) string {
	var username string
	err := db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		fmt.Println("Error getting username:", err)
		return ""
	}
	return username
}