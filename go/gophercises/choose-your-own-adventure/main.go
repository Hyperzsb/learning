package main

import "log"

type Section struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func main() {
	if err := HTMLDemo(); err != nil {
		log.Fatal(err)
	}
}
