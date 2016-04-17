package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/unixpickle/kahoot-hack/kahoot"
)

var token string

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: html <game pin> <nickname>")
		os.Exit(1)
	}
	gamePin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid game pin:", os.Args[1])
		os.Exit(1)
	}
	nickname := os.Args[2]

	http, err := http.Get("https://kahoot.it/reserve/session/"+strconv.Itoa(gamePin))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	token = http.Header.Get("X-Kahoot-Session-Token")

	for _, prefix := range []string{"<h1>", "<u>", "<h2>", "<marquee>", "<button>",
		"<input>", "<pre>", "<textarea>"} {
		if conn, err := kahoot.NewConn(gamePin, token); err != nil {
			fmt.Fprintln(os.Stderr, "failed to connect:", err)
			os.Exit(1)
		} else {
			defer conn.GracefulClose()
			conn.Login(prefix + nickname)
		}
	}

	fmt.Println("Kill this process to deauthenticate.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
