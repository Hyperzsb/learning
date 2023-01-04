package main

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func traverse(root *html.Node, baseUrl *url.URL, linkMap *map[string]bool) error {
	if root.Type == html.ElementNode && root.Data == "a" {
		for _, attr := range root.Attr {
			if attr.Key == "href" {
				newUrl, err := url.Parse(attr.Val)
				if err != nil {
					return err
				}

				if len(newUrl.Scheme) != 0 && newUrl.Scheme != "http" && newUrl.Scheme != "https" {
					break
				}

				if len(newUrl.Host) != 0 && newUrl.Host != baseUrl.Host {
					break
				}

				if len(newUrl.Fragment) != 0 {
					newUrl.Fragment = ""
				}

				link := ""
				newUrl.Scheme = baseUrl.Scheme
				newUrl.Host = baseUrl.Host
				link = newUrl.String()

				if link[len(link)-1] != '/' {
					(*linkMap)[link+"/"] = true
				} else {
					(*linkMap)[link] = true
				}

				break
			}
		}
	}

	for next := root.FirstChild; next != nil; next = next.NextSibling {
		err := traverse(next, baseUrl, linkMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func bfs(targetUrl string, baseUrl *url.URL, linkMap *map[string]bool) error {
	fmt.Printf("Searching in %s\n", targetUrl)

	response, err := http.Get(targetUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		return err
	}

	newLinkMap := make(map[string]bool)
	err = traverse(doc, baseUrl, &newLinkMap)
	if err != nil {
		return err
	}

	for link, _ := range newLinkMap {
		if !(*linkMap)[link] {
			(*linkMap)[link] = true
			err = bfs(link, baseUrl, linkMap)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func SitemapBuilder() error {
	const (
		targetUrl      = "https://hyperzsb.io/"
		outputFilename = ".data/sitemap.xml"
		xmlNs          = "https://www.sitemaps.org/schemas/sitemap/0.9"
	)

	baseUrl, err := url.Parse(targetUrl)
	if err != nil {
		return err
	}

	linkMap := make(map[string]bool)
	linkMap[baseUrl.String()] = true

	err = bfs(targetUrl, baseUrl, &linkMap)
	if err != nil {
		return err
	}

	urlSet := UrlSet{XMLNs: xmlNs, Urls: make([]Url, 0)}
	for link, _ := range linkMap {
		urlSet.Urls = append(urlSet.Urls, Url{Loc: link})
	}

	_ = os.Remove(outputFilename)
	outputFile, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)

	_, err = io.WriteString(outputFile, xml.Header)
	if err != nil {
		return err
	}

	encoder := xml.NewEncoder(outputFile)
	encoder.Indent("", "  ")
	err = encoder.Encode(&urlSet)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := SitemapBuilder(); err != nil {
		log.Fatal(err)
	}
}
