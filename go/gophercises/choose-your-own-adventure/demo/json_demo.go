package demo

import (
	"encoding/json"
	"fmt"
	"os"
)

type Section struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func JSONDemo() error {
	const (
		filename string = ".data/story.json"
	)

	var (
		story = make(map[string]Section)
	)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&story); err != nil {
		return err
	}

	for title, section := range story {
		fmt.Println(title, section.Title)
		for _, option := range section.Options {
			fmt.Println(option.Arc)
		}
		fmt.Println()
	}

	return nil
}
