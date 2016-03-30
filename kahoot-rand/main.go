package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/unixpickle/kahoot-hack/kahoot"
)

var wg sync.WaitGroup

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
	delaystr := "100ms"

	if len(os.Args) == 3 {
		contents, err := ioutil.ReadFile(os.Args[2])
	        if err != nil {
	                fmt.Fprintln(os.Stderr, err)
	                os.Exit(1)
	        }
		res := strings.Split(string(contents), "\n")
		for i := 0; i < len(res); i++ {
			res[i] = strings.TrimSpace(res[i])
			if len(res[i]) == 0 {
				res[i] = res[len(res)-1]
				res = res[:len(res)-1]
				i--
			}
			d, err := time.ParseDuration(delaystr)
			if err != nil {
				fmt.Fprintln(os.Stderr, "invalid delay string:", delaystr)
				fmt.Fprintln(os.Stderr, "valid delay strings include: 250ms, 1s")
				os.Exit(1)
			}
			time.Sleep(d)
			wg.Add(1)
			go launchConnection(gamePin, res[i])
		}
	} else {
		count, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid count:", os.Args[3])
			os.Exit(1)
		}
		nickname := os.Args[2]
		for i := 0; i < count; i++ {
			wg.Add(1)
			d, err := time.ParseDuration(delaystr)
			if err != nil {
				fmt.Fprintln(os.Stderr, "invalid delay string:", delaystr)
				fmt.Fprintln(os.Stderr, "valid delay strings include: 100ms, 0.1s")
			os.Exit(1)
	                }
	                time.Sleep(d)
			go launchConnection(gamePin, nickname+strconv.Itoa(i+1))
		}
	}

	fmt.Println("Terminate this program to stop the automatons...")
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
			quiz.Send(rand.Intn(action.NumAnswers))
		}
	}
}
