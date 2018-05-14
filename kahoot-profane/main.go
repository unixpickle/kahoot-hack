package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/unixpickle/kahoot-hack/kahoot"
)

const ConcurrencyCount = 4

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage: profane <game pin> <nickname prefix> <count>")
		os.Exit(1)
	}

	gamePin := os.Args[1]

	var dieLock sync.Mutex
	connChan := make(chan *kahoot.Conn)
	for i := 0; i < ConcurrencyCount; i++ {
		go func() {
			for {
				conn, err := kahoot.NewConn(gamePin)
				if err != nil {
					dieLock.Lock()
					fmt.Fprintln(os.Stderr, "failed to connect:", err)
					os.Exit(1)
					dieLock.Unlock()
				}
				connChan <- conn
			}
		}()
	}

	for _, nickname := range nicknames() {
		conn := <-connChan
		defer conn.GracefulClose()
		conn.Login(nickname)
	}

	fmt.Println("Kill this process to deauthenticate.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

func nicknames() []string {
	count, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid count:", os.Args[3])
		os.Exit(1)
	}
	base, err := sanitizeName(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to sanitize nickname:", err)
		os.Exit(1)
	}
	res := make([]string, count)
	for x := 0; x < count; x++ {
		res[x] = base + strconv.Itoa(x+1)
	}
	return res
}

func sanitizeName(name string) (string, error) {
	for interval := 16; interval > 1; interval-- {
		attempt := boldify(name, interval)
		testURL := "https://play.kahoot.it/profanity/profanities/nickname/" +
			url.QueryEscape(attempt)
		resp, err := http.Get(testURL)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			return "", err
		}
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		if string(contents) == "[]" {
			return attempt, nil
		}
	}
	return boldify(name, 1), nil
}

func boldify(name string, interval int) string {
	var res []rune
	for i, x := range name {
		if i%interval != 0 {
			res = append(res, x)
			continue
		}
		switch true {
		case x >= 'a' && x <= 'z':
			res = append(res, 120302+x-'a')
		case x >= 'A' && x <= 'Z':
			res = append(res, 120276+x-'A')
		default:
			res = append(res, x)
		}
	}
	return string(res)
}
