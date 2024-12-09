# goblet

A simple and lightweight web framework for Go, that works like Flask, but in an ambient that is much faster than Python.

## Installation

```bash
go get github.com/aquiffoo/goblet
```

## Usage

```go
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
		}

		app.Render(w, "index.html", data)
	})

	app.Handle("/about", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"title": "About",
		}

		app.Render(w, "about.html", data)
	})

	err := app.Serve("8080")
	if err != nil {
		panic(err)
	}
}
```

## License

This project is licensed under the [CC BY-NC 4.0](./LICENSE) license.

## Authorship
This was entirely written by [aquiffoo](https://github.com/aquiffoo). Contributors can add their names to the [AUTHORS](./AUTHORS) file.
