package handler

import (
	"fmt"
	"net/http"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) (http.HandlerFunc, error) {
	if fallback == nil {
		return nil, fmt.Errorf("fallback handler not specified")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if target, exist := pathsToUrls[r.URL.Path]; exist {
			http.Redirect(w, r, target, 303)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}, nil
}
