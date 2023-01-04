package demo

import (
	"encoding/xml"
	"fmt"
	"os"
)

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Urls    []Url    `xml:"url"`
}

type Url struct {
	Loc string `xml:"loc"`
}

func XMLDemo() error {
	const (
		inputFilename  string = ".data/input.xml"
		outputFilename string = ".data/output.xml"
	)

	var (
		siteMap UrlSet
	)

	// Decoder demo
	inputFile, err := os.Open(inputFilename)
	if err != nil {
		return err
	}

	decoder := xml.NewDecoder(inputFile)
	err = decoder.Decode(&siteMap)
	if err != nil {
		return err
	}

	fmt.Println(siteMap)

	// Encoder demo
	outputFile, err := os.OpenFile(outputFilename, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err

	}

	encoder := xml.NewEncoder(outputFile)
	encoder.Indent("", "    ")
	err = encoder.Encode(&siteMap)
	if err != nil {
		return err
	}

	return nil
}
