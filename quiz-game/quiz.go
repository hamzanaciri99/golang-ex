package main

import (
	"fmt"
	"flag"
	"os"
	"bufio"
	"math/rand"
	"time"
	"io"
	"strings"
	"github.com/hamzanaciri99/goland-ex/util"
)

//Define and initialize flags
var (
	shuffle bool

)
func flagsInit() {
	flag.BoolVar(&shuffle, "shuffle", false, "set to true to shuffle the quiz order")
}

func init() {
	flagsInit()
}	

type Quiz struct {
	expr, val string
}

func Suffle[T interface{}](q []T) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(q), func(i, j int) { q[i], q[j] = q[j], q[i] })
}

func main() {
	flag.Parse()

	s, err := os.ReadFile("/Users/azmah/Desktop/hamzanaciri99/golang-exercices/quiz-game/quiz.csv")
	util.CheckError(err)

	lines := strings.Split(string(s), "\n")

	result := 0
	count := 0
	quizzes := []Quiz{}

	for _, line := range lines {
		split := strings.Split(line, ",")
		if len(split) != 2 {
			continue
		}
		quizzes = append(quizzes, Quiz{split[0], split[1]})
	}

	if shuffle {
		Suffle(quizzes)
	}

	for _, quiz := range quizzes {
		fmt.Printf("What is %s?\n", quiz.expr)
		
		stdin := bufio.NewReader(os.Stdin)
		val, err := stdin.ReadString('\n')
		util.CheckError(err, io.EOF)

		count++
		if strings.TrimSpace(val) == strings.TrimSpace(quiz.val) {
			result++
		} else {
			break
		}
	}

	fmt.Printf("Your score is %d out of %d\n", result, count)

}