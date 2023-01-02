package demo

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Mapping struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func YAMLDemo() error {
	const (
		filename string = ".data/mappings.yaml"
	)
	var (
		mappings []Mapping
	)

	// TODO: figure out the differences between os.ReadFile and os.Open,
	//       as well as the concept and usage of io.Reader

	// Use os.ReadFile and yaml.Unmarshal to read and parse yaml file
	fileBuff, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileBuff, &mappings)
	if err != nil {
		return err
	}

	for _, mapping := range mappings {
		fmt.Println(mapping.Path, mapping.Url)
	}

	// Use os.Open and yaml.Decoder to read and parse yaml file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&mappings)
	if err != nil {
		return err
	}

	for _, mapping := range mappings {
		fmt.Println(mapping.Path, mapping.Url)
	}

	return nil
}
