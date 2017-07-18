package main

//written by Peter Stenger on July 17, 2017 (@reteps)
import (
	"fmt"
	"github.com/unixpickle/kahoot-hack/kahoot"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) != 6 {
		fmt.Fprintln(os.Stderr, "Usage: auto <quizid> <game pin> <nickname> <email> <password>")
		os.Exit(1)
	}

	gamePin, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid game pin:", os.Args[1])
		os.Exit(1)
	}
	nickname := os.Args[3]
	quizid := os.Args[1]
	email := os.Args[4]
	password := os.Args[5]
	//get access token
	token := kahoot.AccessToken(email, password)
	//get all information from quiz
	data := kahoot.ReturnData(token, quizid)
	//return question data from quiz
	answers := kahoot.ParseData(data)

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
	fmt.Println("waiting to start...")
	questionnum := 0
	for {

		action, err := quiz.Receive()
		if err != nil {
			if !<-closed {
				fmt.Fprintln(os.Stderr, "Could not receive question:", err)
			}
			os.Exit(1)
		}
		if action.Type == kahoot.QuestionIntro {
			fmt.Printf("Question %d starting...\n", questionnum+1)
		} else if action.Type == kahoot.QuestionAnswers {
			answer, _ := strconv.Atoi(answers[questionnum][2])
			if err := quiz.Send(answer); err != nil {
				fmt.Fprintln(os.Stderr, "Could not answer:", err)
				os.Exit(1)
			}
			fmt.Printf("Answered %s (%s)\n", answers[questionnum][1], answers[questionnum][3])
			questionnum += 1
		}
	}
}
