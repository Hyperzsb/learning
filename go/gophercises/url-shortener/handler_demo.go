package main

import (
	"gophercises/urlshortener/handler"
	"log"
	"net/http"
)

func DemoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/google" {
			http.Redirect(w, r, "https://google.com", 303)
		} else {
			http.NotFound(w, r)
		}
	}
}

func main() {
	pathsToUrls := map[string]string{
		"/google": "https://google.com/",
		"/go":     "https://go.dev",
	}

	mapHandler, err := handler.MapHandler(pathsToUrls, http.NotFoundHandler())
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8080", mapHandler); err != nil {
		log.Fatal(err)
	}
}
