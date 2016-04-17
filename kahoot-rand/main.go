package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/unixpickle/kahoot-hack/kahoot"
)

var wg sync.WaitGroup

const ConnectionDelay = time.Millisecond * 100

var StatisticsChan = make(chan int, 0)
var AnswerCount uint32

func main() {
	if len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage: rand <game pin> <nickname prefix> <count>")
		fmt.Fprintln(os.Stderr, "       rand <game pin> <name_list.txt>")
		os.Exit(1)
	}

	gamePin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid game pin:", os.Args[1])
		os.Exit(1)
	}

	nicknames, err := readNicknames()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	botCount := 0
	for nickname := range nicknames {
		wg.Add(1)
		go launchConnection(gamePin, nickname)
		time.Sleep(ConnectionDelay)
		botCount++
	}

	fmt.Println("Entered with", botCount, "bots.")
	fmt.Println("Terminate this program to stop the automatons...")

	go printStatistics(botCount)

	wg.Wait()
}

func launchConnection(gamePin int, nickname string) {
	defer wg.Done()

	conn, err := kahoot.NewConn(gamePin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to connect:", err)
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

	if err := conn.Login(nickname); err != nil {
		fmt.Fprintln(os.Stderr, "failed to login:", err)
		os.Exit(1)
	}

	quiz := kahoot.NewQuiz(conn)

	for {
		action, err := quiz.Receive()
		if err != nil {
			if <-closed {
				return
			} else {
				fmt.Fprintln(os.Stderr, "Could not receive question:", err)
				os.Exit(1)
			}
		}
		if action.Type == kahoot.QuestionAnswers {
			atomic.StoreUint32(&AnswerCount, uint32(action.NumAnswers))
			answer := rand.Intn(action.NumAnswers)
			quiz.Send(action.AnswerMap[answer])
			StatisticsChan <- answer
		}
	}
}

func printStatistics(clientCount int) {
	qid := 1
	for {
		stats := map[int]int{}
		for i := 0; i < clientCount; i++ {
			stats[<-StatisticsChan]++
		}

		answerCount := int(atomic.LoadUint32(&AnswerCount))

		fmt.Println("-----STATISTICS-----")
		fmt.Println("Question ID:", qid)
		for i := 0; i < answerCount; i++ {
			fmt.Println("Answer "+strconv.Itoa(i)+":", stats[i])
		}
		qid++
	}
}

func readNicknames() (<-chan string, error) {
	if len(os.Args) == 3 {
		contents, err := ioutil.ReadFile(os.Args[2])
		if err != nil {
			return nil, err
		}
		nameLines := strings.Split(string(contents), "\n")
		res := make(chan string, len(nameLines))
		for _, line := range nameLines {
			nickname := strings.TrimSpace(line)
			if len(nickname) != 0 {
				res <- nickname
			}
		}
		close(res)
		return res, nil
	}

	count, err := strconv.Atoi(os.Args[3])
	if err != nil {
		return nil, errors.New("invalid count: " + os.Args[3])
	}

	baseName := os.Args[2]
	res := make(chan string)
	go func() {
		for i := 0; i < count; i++ {
			res <- baseName + strconv.Itoa(i+1)
		}
		close(res)
	}()
	return res, nil
}
