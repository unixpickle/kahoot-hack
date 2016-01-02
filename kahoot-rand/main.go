package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/unixpickle/kahoot-hack/kahoot"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage: rand <game pin> <nickname prefix> <count>")
		os.Exit(1)
	}
	gamePin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid game pin:", os.Args[1])
		os.Exit(1)
	}

	count, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid count:", os.Args[3])
		os.Exit(1)
	}

	nickname := os.Args[2]

	for i := 0; i < count; i++ {
		go launchConnection(gamePin, nickname+strconv.Itoa(i+1))
	}

	fmt.Println("Terminate this program to stop the automatons...")
	select {}
}

func launchConnection(gamePin int, nickname string) {
	conn, err := kahoot.NewConn(gamePin)
	defer conn.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to connect:", err)
		os.Exit(1)
	}
	if err := conn.Login(nickname); err != nil {
		fmt.Fprintln(os.Stderr, "failed to login:", err)
		os.Exit(1)
	}
	quiz := kahoot.NewQuiz(conn)
	for {
		action, err := quiz.Receive()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not receive question:", err)
			os.Exit(1)
		}
		if action.Type == kahoot.QuestionAnswers {
			quiz.Send(rand.Intn(action.NumAnswers))
		}
	}
}
