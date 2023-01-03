package demo

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func traverse(root *html.Node) {
	if root.Type == html.ElementNode {
		fmt.Println("ElementNode: ", root.Data)
	}
	if root.Type == html.TextNode {
		if len(strings.Trim(root.Data, " \n")) != 0 {
			fmt.Println("TextNode: ", root.Data)
		}
	}
	for next := root.FirstChild; next != nil; next = next.NextSibling {
		traverse(next)
	}
}

func HTMLDemo() error {
	const (
		filename = ".data/example-1.html"
	)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	doc, err := html.Parse(file)
	if err != nil {
		return err
	}

	traverse(doc)

	return nil
}
