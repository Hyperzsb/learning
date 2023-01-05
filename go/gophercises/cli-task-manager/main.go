package main

import (
	"gophercises/taskmanager/demo"
	"log"
)

func main() {
	if err := demo.CobraDemo(); err != nil {
		log.Fatal(err)
	}
}
