package main

import (
	"gophercises/taskmanager/demo"
	"log"
)

func main() {
	if err := demo.BBoltDemo(); err != nil {
		log.Fatal(err)
	}
}
