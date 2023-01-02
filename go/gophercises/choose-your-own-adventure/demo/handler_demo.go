package main

import (
	"log"
	"net/http"
)

type MyHandler struct {
	HTML string
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(h.HTML)); err != nil {
		log.Fatal(err)
	}
}

func HandlerDemo() error {
	handler := MyHandler{HTML: "<h1>Hello World</h1>"}
	err := http.ListenAndServe(":8080", &handler)
	if err != nil {
		return err
	}

	return nil
}
