package main

import (
	"log"
	"net/http"
)

func Serve() error {
	const (
		storyJSONFilename     string = ".data/story.json"
		storyTemplateFilename string = "template.html"
		servePort             string = ":8080"
	)
	storyContent, err := NewStory(storyJSONFilename)
	if err != nil {
		return err
	}

	storyTemplate, err := NewStoryTemplate(storyTemplateFilename)
	if err != nil {
		return err
	}

	storyHandler, err := NewStoryHandler(storyContent, storyTemplate)
	if err != nil {
		return err
	}

	err = http.ListenAndServe(servePort, storyHandler)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Serve(); err != nil {
		log.Fatal(err)
	}
}
