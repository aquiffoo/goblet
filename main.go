package goblet

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Goblet struct {
	routes    map[string]http.HandlerFunc
	templates *template.Template
	static    string
	hotReload bool
}

func loadTemplates(path string) *template.Template {
	funcMap := template.FuncMap{
		"url_for": UrlFor,
		"block":    func(string, interface{}) string { return "" },
		"define":   func(string) string { return "" },
		"if":       func(interface{}) string { return "" },
		"else":     func() string { return "" },
		"end":      func() string { return "" },
		"range":    func(interface{}) string { return "" },
		"with":     func(interface{}) string { return "" },
		"template": func(string, ...interface{}) string { return "" },
		"extends":  func(string, ...interface{}) string { return "" },
	}
	p := filepath.Join(path, "*.html")
	return template.Must(template.New("").Funcs(funcMap).ParseGlob(p))
}

func New(hotReload bool) *Goblet {
	g := &Goblet{
		routes:    make(map[string]http.HandlerFunc),
		templates: loadTemplates("templates"),
		static:    "./static",
		hotReload: hotReload,
	}

	if hotReload {
		go g.watch()
	}

	return g
}

func (g *Goblet) Handle(path string, handler http.HandlerFunc) {
	g.routes[path] = handler
}

func (g *Goblet) Serve(port string) error {
	mux := http.NewServeMux()
	for path, handler := range g.routes {
		mux.HandleFunc(path, handler)
	}

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(g.static))))

	fmt.Printf("Goblet serving at %s\n", port)
	return http.ListenAndServe(":"+port, mux)
}

func (g *Goblet) Render(w http.ResponseWriter, name string, data interface{}) error {
    w.Header().Set("Content-Type", "text/html")

    finalTmpl, err := Extends(g.templates, name, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return err
    }

    err = finalTmpl.ExecuteTemplate(w, name, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return err
    }

    return nil
}

func (g *Goblet) watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERR: Failed to create watcher:", err)
		return
	}
	defer watcher.Close()

	watchPaths := []string{"."}
	excludedDirs := map[string]bool{".git": true, "bin": true, "node_modules": true}
	excludedFiles := map[string]bool{"main.exe": true}
	fileExtensions := map[string]bool{".go": true, ".html": true, ".css": true, ".js": true}

	for _, path := range watchPaths {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				if excludedDirs[info.Name()] {
					return filepath.SkipDir
				}
				err := watcher.Add(path)
				if err != nil {
					fmt.Printf("ERR:  Failed to watch directory %s: %v\n", path, err)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("ERR:  Failed to walk path %s: %v\n", path, err)
		}
	}

	fmt.Println("GOBLET: Hot reload enabled. Watching for changes...")

	restartTimer := time.NewTimer(0)
	restartTimer.Stop()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				fileName := filepath.Base(event.Name)
				ext := filepath.Ext(event.Name)

				if excludedFiles[fileName] || !fileExtensions[ext] {
					continue
				}

				fmt.Printf("Change detected in %s\n", event.Name)

				restartTimer.Reset(10 * time.Millisecond)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("ERR:  Watcher error:", err)

		case <-restartTimer.C:
			fmt.Println("Restarting server...")

			cmd := exec.Command("go", "run", "main.go")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Start()
			if err != nil {
				fmt.Printf("ERR:  Error restarting server: %v\n", err)
			}

			os.Exit(0)
		}
	}
}