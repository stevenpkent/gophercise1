package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	fileName        string
	timeLimit       int
	needToRandomize bool
)

func init() {
	flag.StringVar(&fileName, "csv", "problems.csv", "question file")
	flag.IntVar(&timeLimit, "timelimit", 30, "amount of time to take the quiz")
	flag.BoolVar(&needToRandomize, "shuffle", false, "jumble questions?")
	flag.Parse()
}

func main() {
	file := getFile(fileName)
	lines := getLines(file)
	questions := getQuestions(lines)
	createQuiz(questions, timeLimit, needToRandomize)
}

func getFile(name string) *os.File {
	file, err := os.Open(name)

	if err == nil {
		exit(fmt.Sprintf("Failed to open the file: %s\n", name))
	}

	return file
}

func getLines(file *os.File) [][]string {
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()

	if err != nil {
		exit("Failed to parse the CSV file")
	}

	return lines
}

func getQuestions(lines [][]string) []question {
	questions := make([]question, len(lines))

	for index, line := range lines {
		questions[index] = question{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return questions
}

func createQuiz(questions []question, timelimit int, needToRandomize bool) {
	numberOfQuestions := len(questions)
	indices := getIndices(needToRandomize, numberOfQuestions)
	var numberOfCorrectAnswers int
	var answer string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Press enter to start quiz")
	reader.ReadString('\n') // start timer when user clicks return
	timer := time.NewTimer(time.Duration(timelimit) * time.Second)

	for index, i := range indices {
		fmt.Printf("Question #%d: %s = ", index+1, questions[i].q)
		answerChannel := make(chan string)

		go func() {
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		case receivedAnswer := <-answerChannel:
			if strings.ToLower(receivedAnswer) == strings.ToLower(questions[i].a) {
				numberOfCorrectAnswers++
			}

		case <-timer.C: // timeout occurred
			fmt.Println()
			outputQuizResult(numberOfCorrectAnswers, numberOfQuestions)
			return // end the quiz
		}
	}

	outputQuizResult(numberOfCorrectAnswers, numberOfQuestions)
}

func getIndices(mustRandomize bool, numberOfQuestions int) []int {
	rand.Seed(time.Now().UTC().UnixNano())
	indices := rand.Perm(numberOfQuestions)

	if !mustRandomize {
		sort.Ints(indices[:])
	}

	return indices
}

func outputQuizResult(numberOfCorrectAnswers int, numberOfQuestions int) {
	fmt.Printf("You scored %d out of %d\n", numberOfCorrectAnswers, numberOfQuestions)
}

func exit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

type question struct {
	q string
	a string
}
