package main

import (
	"net/http"

	"github.com/aquiffoo/goblet"
)

func main() {
	app := goblet.New(true)

	app.Handle("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"title": "Home",
			"lang": "en",
			"charset": "utf-8",
		}

		app.Render(w, "index.html", data)
	})

	app.Handle("/about", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"title": "Home",
			"lang": "en",
			"charset": "utf-8",
		}

		app.Render(w, "about.html", data)
	})

	err := app.Serve("8080")
	if err != nil {
		panic(err)
	}
}
