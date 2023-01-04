package main

import (
	"gophercises/sitemapbuilder/demo"
	"log"
)

func main() {
	if err := demo.XMLDemo(); err != nil {
		log.Fatal(err)
	}
}
