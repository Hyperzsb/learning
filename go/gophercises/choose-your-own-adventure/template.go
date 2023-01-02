package main

import (
	"html/template"
	"io"
	"os"
)

func NewStoryTemplate(filename string) (*template.Template, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	tmplStr, tmplChars := "", make([]byte, 512)
	for {
		if _, err := file.Read(tmplChars); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else {
			tmplStr += string(tmplChars)
		}
	}

	tmpl, err := template.New("StoryTemplate").Parse(tmplStr)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
