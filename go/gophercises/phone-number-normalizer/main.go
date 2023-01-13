package main

import (
	"gophercises/numbernormalizer/demo"
	"log"
)

func main() {
	if err := demo.GORMDemo(); err != nil {
		log.Fatal(err)
	}
}
