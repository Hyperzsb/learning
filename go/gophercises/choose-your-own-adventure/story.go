package main

import (
	"encoding/json"
	"os"
)

type Story map[string]Section

type Section struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func NewStory(filename string) (*Story, error) {
	var story Story

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&story)
	if err != nil {
		return nil, err
	}

	return &story, nil
}
