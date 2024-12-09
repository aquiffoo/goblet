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

```html
<!-- index.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Home</title>
    <link rel="stylesheet" href="{{ url_for "static/style.css" }}">
</head>
<body>
    <h1>Welcome to Goblet</h1>
    <p>This is a demo for Goblet, my first Go framework.</p>
    <a href="{{ url_for "about" }}">About</a>
</body>
</html>

```

```html
<!-- about.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>About</title>
    <link rel="stylesheet" href="{{ url_for "static/style.css" }}">
</head>
<body>
    <h1>About Page</h1>
    <p>This page is rendered using Goblet templates.</p>
    <a href="{{ url_for "index" }}">Go Back</a>
</body>
</html>

```

## Documentation

The documentation can be found [here](https://godoc.org/github.com/aquiffoo/goblet).

## License

This project is licensed under the [CC BY-NC 4.0](./LICENSE) license.

## Authorship
This was entirely written by [aquiffoo](https://github.com/aquiffoo). Contributors can add their names to the [AUTHORS](./AUTHORS) file.
