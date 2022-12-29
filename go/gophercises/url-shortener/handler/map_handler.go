package handler

import (
	"net/http"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) (http.HandlerFunc, error) {
	return nil, nil
}
