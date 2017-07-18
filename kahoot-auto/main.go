package main

//written by Peter Stenger on July 17, 2017 (@reteps)
import (
	"fmt"
	"kahoot"
	//"github.com/unixpickle/kahoot-hack/kahoot"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// ParseQuizInformation parses quiz information
// from a kahoot. It returns the question, answer,
// default answer number, and default answer color.
func ParseQuizInformation(data *kahoot.KahootQuiz) [][]string {
	var results [][]string
	colormap := map[int]string{0: "red", 1: "blue", 2: "yellow", 3: "blue"}
	for _, value := range data.Questions {
		var questiondata []string
		for i, choice := range value.Choices {
			if choice.Correct == true {
				questiondata = append(questiondata, value.Question, choice.Answer, strconv.Itoa(i), colormap[i])
				break
			}
		}
		results = append(results, questiondata)
	}
	return results
}
func Prompt(question string) string {
	fmt.Print(question)
	var response string
	fmt.Scanf("%s", &response)
	return response
}
func main() {
	argnum := len(os.Args)
	if argnum != 5 && argnum != 4 {
		fmt.Fprintln(os.Stderr, "Usage: auto <quizid> <game pin> <nickname> (email)")
		os.Exit(1)
	}

	gamePin, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid game pin:", os.Args[1])
		os.Exit(1)
	}
	nickname := os.Args[3]
	quizid := os.Args[1]
	var email string
	if argnum == 4 {
		email = Prompt("email > ")
	} else {
		email = os.Args[4]
	}
	password := Prompt("password > ")
	token, err := kahoot.AccessToken(email, password)
	if err != nil {
		panic(err)
	}
	data := kahoot.QuizInformation(token, quizid)
	answers := ParseQuizInformation(data)

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
