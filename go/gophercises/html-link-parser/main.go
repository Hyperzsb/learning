package main

import (
	"golang.org/x/net/html"
	"log"
	"os"
	"strings"
)

func traverse(root *html.Node, linSet *LinkSet) {
	if root.Type == html.ElementNode && root.Data == "a" {
		href, text := "", ""
		for _, attr := range root.Attr {
			if attr.Key == "href" {
				href = attr.Val
				break
			}
		}
		text = getText(root)
		linSet.Add(href, text)
	}

	for next := root.FirstChild; next != nil; next = next.NextSibling {
		traverse(next, linSet)
	}
}

func getText(root *html.Node) string {
	text := ""

	for next := root.FirstChild; next != nil; next = next.NextSibling {
		text += getText(next)
	}

	if root.Type == html.TextNode && len(strings.Trim(root.Data, " \n")) != 0 {
		text = strings.Trim(root.Data, " \n") + text
	}

	return text
}

func LinkParser() error {
	const (
		filename string = ".data/example-4.html"
	)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	doc, err := html.Parse(file)
	if err != nil {
		return err
	}

	linkSet := NewLinkSet()

	traverse(doc, linkSet)

	linkSet.Print()

	return nil
}

func main() {
	if err := LinkParser(); err != nil {
		log.Fatal(err)
	}
}
