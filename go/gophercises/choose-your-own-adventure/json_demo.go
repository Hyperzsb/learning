package main

import (
	"encoding/json"
	"fmt"
	"os"
)

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
