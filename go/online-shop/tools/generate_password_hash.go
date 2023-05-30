package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func main() {
	rawPassword := ""
	_, _ = fmt.Scanf("%s", &rawPassword)

	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(hash))

	err = bcrypt.CompareHashAndPassword(hash, []byte(rawPassword))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Matched")
}
