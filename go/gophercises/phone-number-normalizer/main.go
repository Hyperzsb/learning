package main

import (
	"gophercises/numbernormalizer/demo"
	"log"
)

func main() {
	if err := demo.PostgresDemo(); err != nil {
		log.Fatal(err)
	}
}
