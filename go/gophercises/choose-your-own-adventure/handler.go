package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type StoryHandler struct {
	Story          *Story
	Template       *template.Template
	initialSection string
}

func (h *StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		err := h.Template.Execute(w, (*(h.Story))[h.initialSection])
		if err != nil {
			http.Error(w, err.Error(), 502)
		}
		return
	}

	if section, exist := (*(h.Story))[(r.URL.Path)[1:]]; exist {
		err := h.Template.Execute(w, section)
		if err != nil {
			http.Error(w, err.Error(), 502)
		}
	} else {
		http.NotFound(w, r)
	}
}

func NewStoryHandler(storyContent *Story, storyTemplate *template.Template) (http.Handler, error) {
	sectionInDegree := make(map[string]int)
	for arc, _ := range *storyContent {
		sectionInDegree[arc] = 0
	}
	for _, section := range *storyContent {
		for _, option := range section.Options {
			if _, exist := sectionInDegree[option.Arc]; !exist {
				return nil, fmt.Errorf("invalid option: %v", option)
			} else {
				sectionInDegree[option.Arc]++
			}
		}
	}
	initialSection := ""
	for section, count := range sectionInDegree {
		if count == 0 {
			if initialSection != "" {
				return nil, fmt.Errorf("duplicate initial section: %v and %v", initialSection, section)
			} else {
				initialSection = section
			}
		}
	}
	if initialSection == "" {
		return nil, fmt.Errorf("no valid initial section")
	}

	if storyContent == nil {
		return nil, fmt.Errorf("story content not provided")
	}
	if storyTemplate == nil {
		return nil, fmt.Errorf("stroy template not specified")
	}

	storyHandler := StoryHandler{Story: storyContent, Template: storyTemplate, initialSection: initialSection}
	return &storyHandler, nil
}
