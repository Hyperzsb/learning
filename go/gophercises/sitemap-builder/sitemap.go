package main

import "encoding/xml"

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNs   string   `xml:"xmlns,attr"`
	Urls    []Url    `xml:"url"`
}

type Url struct {
	Loc string `xml:"loc"`
}
