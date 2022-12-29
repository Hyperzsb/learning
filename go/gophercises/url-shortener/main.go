package main

import (
	"fmt"
	"gophercises/urlshortener/handler"
	"log"
	"net/http"
)

func UrlShortener() error {
	// Default fallback
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello, world!")
	})

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/fmt": "https://pkg.go.dev/fmt",
		"/log": "https://pkg.go.dev/log",
	}
	mapHandler, err := handler.MapHandler(pathsToUrls, mux)
	if err != nil {
		return err
	}

	// Build the YAMLHandler using the mapHandler as the fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := handler.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		return err
	}

	err = http.ListenAndServe(":8080", yamlHandler)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := UrlShortener(); err != nil {
		log.Fatal(err)
	}
}
