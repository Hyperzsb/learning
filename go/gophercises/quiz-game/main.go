package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"gophercises/quizgame/question"
	"log"
	"os"
	"time"
)

func Quiz(filename string, questionNum, correctNum *int) error {
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvRecords, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	questions := make([]question.Question, 0, len(csvRecords))
	for idx, record := range csvRecords {
		questions = append(questions, question.Question{})
		if err := questions[idx].Parse(record); err != nil {
			return err
		}
	}
	*questionNum = len(questions)

	for i, q := range questions {
		// Get user input and validate it
		userAns := 0
		fmt.Printf("Q #%d. %d%s%d = ", i+1, q.Operand[0], string(q.Operator), q.Operand[1])
		_, err = fmt.Scanf("%d", &userAns)
		for err != nil {
			fmt.Printf("Please provide a valid numeric answer: %s\n", err)
			fmt.Printf("Q #%d. %d%s%d = ", i+1, q.Operand[0], string(q.Operator), q.Operand[1])
			_, err = fmt.Scanf("%d", &userAns)
		}

		if userAns == q.Result {
			*correctNum++
		}
	}

	return nil
}

func QuizGame() error {
	// Define cli flags
	var (
		filename string
	)
	// Parse cli flags
	flag.StringVar(&filename, "file", ".data/questions.csv", "question file in .csv format")
	flag.Parse()

	questionNum, correctNum := 0, 0
	if err := Quiz(filename, &questionNum, &correctNum); err != nil {
		return err
	} else {
		fmt.Printf("\nQuiz completed! Your score: %d out of %d\n", correctNum, questionNum)
		return nil
	}
}

func QuizGameWithTimer() error {
	// Define cli flags
	var (
		filename  string
		timeLimit time.Duration
	)
	// Parse cli flags
	flag.StringVar(&filename, "file", ".data/questions.csv", "question file in .csv format")
	flag.DurationVar(&timeLimit, "time", 30*time.Second, "time limit for the quiz")
	flag.Parse()

	done := make(chan bool)
	errChan := make(chan error)

	questionNum, correctNum := 0, 0
	go func() {
		if err := Quiz(filename, &questionNum, &correctNum); err != nil {
			errChan <- err
		} else {
			done <- true
		}
	}()

	select {
	case <-done:
		fmt.Printf("\nQuiz completed! Your score: %d out of %d\n", correctNum, questionNum)
		return nil
	case err := <-errChan:
		return err
	case <-time.After(timeLimit):
		fmt.Printf("\nTime up! Your score: %d out of %d\n", correctNum, questionNum)
		return nil
	}
}

func main() {
	if err := QuizGameWithTimer(); err != nil {
		log.Fatal(err)
	}
}
