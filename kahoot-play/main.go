package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/unixpickle/kahoot-hack/kahoot"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: play <game pin> <nickname>")
		os.Exit(1)
	}

	gamePin := os.Args[1]
	nickname := os.Args[2]

	conn, err := kahoot.NewConn(gamePin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to connect:", err)
		os.Exit(1)
	}
	if err := conn.Login(nickname); err != nil {
		fmt.Fprintln(os.Stderr, "failed to login:", err)
		os.Exit(1)
	}

	closed := make(chan bool, 1)
	closed <- false
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		<-closed
		closed <- true
		conn.GracefulClose()
	}()

	quiz := kahoot.NewQuiz(conn)
	for {
		action, err := quiz.Receive()
		if err != nil {
			if !<-closed {
				fmt.Fprintln(os.Stderr, "Could not receive question:", err)
			}
			os.Exit(1)
		}
		if action.Type == kahoot.QuestionIntro {
			fmt.Println("Awaiting answers...")
		} else if action.Type == kahoot.QuestionAnswers {
			fmt.Print("Answer (0 through " + strconv.Itoa(action.NumAnswers-1) + "): ")
			answer := readNumberInput()
			if err := quiz.Send(answer); err != nil {
				fmt.Fprintln(os.Stderr, "Could not answer:", err)
				os.Exit(1)
			}
		}
	}
}

func readNumberInput() int {
	for {
		var buffer string
		for {
			buf := make([]byte, 1)
			if _, err := os.Stdin.Read(buf); err != nil {
				panic("could not read input")
			}
			if buf[0] == '\r' {
				continue
			} else if buf[0] == '\n' {
				break
			}
			buffer += string(rune(buf[0]))
		}
		res, err := strconv.Atoi(buffer)
		if err != nil {
			fmt.Println("please enter a number")
			continue
		} else {
			return res
		}
	}
}
