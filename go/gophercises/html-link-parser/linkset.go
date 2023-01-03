package main

import "fmt"

type LinkSet struct {
	Links []Link
}

type Link struct {
	Href string
	Text string
}

func NewLinkSet() *LinkSet {
	return &LinkSet{make([]Link, 0)}
}

func (ls *LinkSet) Add(href, text string) {
	ls.Links = append(ls.Links, Link{href, text})
}

func (ls *LinkSet) Print() {
	for _, link := range ls.Links {
		fmt.Printf("Link {\n    Href: \"%s\"\n    Text: \"%s\"\n}\n", link.Href, link.Text)
	}
}
