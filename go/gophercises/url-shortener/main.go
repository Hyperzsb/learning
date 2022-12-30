package main

import (
	"fmt"
	"gophercises/urlshortener/handler"
	"log"
	"net/http"
	"os"
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
	// Use raw string as yaml input
	/*
			yaml := `
		- path: /net/http
		  url: https://pkg.go.dev/net/http
		- path: /strings
		  url: https://pkg.go.dev/strings
		`
			yamlHandler, err := handler.YAMLHandler(strings.NewReader(yaml), mapHandler)
			if err != nil {
				return err
			}
	*/
	// Or, use file as yaml input
	yaml, err := os.Open(".data/mappings.yaml")
	if err != nil {
		return err
	}

	yamlHandler, err := handler.YAMLHandler(yaml, mapHandler)
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
