package handler

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
)

type Mapping struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func YAMLHandler(reader io.Reader, fallback http.Handler) (http.HandlerFunc, error) {
	var (
		mappings []Mapping
	)

	if fallback == nil {
		return nil, fmt.Errorf("fallback handler not specified")
	}

	decoder := yaml.NewDecoder(reader)
	err := decoder.Decode(&mappings)
	if err != nil {
		return fallback.ServeHTTP, err
	}

	mappingMap := make(map[string]string)
	for _, mapping := range mappings {
		mappingMap[mapping.Path] = mapping.Url
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if url, exist := mappingMap[r.URL.Path]; exist {
			http.Redirect(w, r, url, 303)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}, nil
}
