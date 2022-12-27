package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"gophercises/quizgame/question"
	"log"
	"os"
	"strings"
)

func QuizGame() error {
	// Define cli flags
	var (
		filename string
	)

	// Parse cli flags
	flag.StringVar(&filename, "file", ".data/questions.csv", "question file in .csv format")
	flag.Parse()

	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse and print questions per line
	q, questionNum, correctNum := question.Question{}, 0, 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		csvReader := csv.NewReader(strings.NewReader(scanner.Text()))
		csvRecord, err := csvReader.Read()
		if err != nil {
			return err
		}

		if err := q.Parse(csvRecord); err != nil {
			return err
		}
		questionNum++

		fmt.Printf("Q #%d. %d%s%d = ", questionNum, q.Operand[0], string(q.Operator), q.Operand[1])

		userAns := 0
		_, err = fmt.Scanf("%d", &userAns)
		for err != nil {
			fmt.Printf("Please provide a valid numeric answer: %s\n", err)
			fmt.Printf("%s=", csvRecord[0])
			_, err = fmt.Scanf("%d", &userAns)
		}

		if userAns == q.Result {
			correctNum++
		}
	}

	fmt.Printf("Your score: %d/%d\n", correctNum, questionNum)

	return nil
}

func main() {
	if err := QuizGame(); err != nil {
		log.Fatal(err)
	}
}
