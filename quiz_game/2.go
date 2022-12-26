package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problems struct {
	question string
	answer   string
}

func main() {
	csvFileName := flag.String("csv", "problems.csv", "Please enter a csv file name, having format question,answer")
	timeLimit := flag.Int("timeLimit", 5, "Give time limit in seconds.")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Error while trying to open: %s\n", *csvFileName))
	}

	r := csv.NewReader(file)
	csvRecords, err := r.ReadAll()

	if err != nil {
		exit(fmt.Sprintf("Error while trying to read CSV records"))
	}

	problemList := parseLines(csvRecords)

	correctCount := 0

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for key, value := range problemList {
		fmt.Printf("%d# %s=", key+1, value.question)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Println("\nSorry your time is up!")
			fmt.Printf("\nYou got %d correct out of %d\n", correctCount, len(problemList))
			return
		case answer := <-answerCh:
			if answer == value.answer {
				correctCount++
			}
		}
	}

	fmt.Printf("\nYou got %d correct out of %d\n", correctCount, len(problemList))
}

func parseLines(csvRecords [][]string) []problems {
	problemList := make([]problems, len(csvRecords))

	for i, row := range csvRecords {
		problemList[i].question = row[0]
		problemList[i].answer = strings.TrimSpace(row[1])
	}

	return problemList
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}
