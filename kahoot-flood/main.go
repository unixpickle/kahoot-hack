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
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage: flood <game pin> <nickname prefix> <count>")
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
		if conn, err := kahoot.NewConn(gamePin); err != nil {
			fmt.Fprintln(os.Stderr, "failed to connect:", err)
			os.Exit(1)
		} else {
			defer conn.GracefulClose()
			conn.Login(nickname + strconv.Itoa(i+1))
		}
	}

	fmt.Println("Kill this process to deauthenticate.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
