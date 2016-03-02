package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/unixpickle/kahoot-hack/kahoot"
)

func main() {
	if len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Usage: flood <game pin> <nickname prefix> <count>")
		fmt.Fprintln(os.Stderr, "       flood <game pin> <name_list.txt>")
		os.Exit(1)
	}
	gamePin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid game pin:", os.Args[1])
		os.Exit(1)
	}

	for _, nickname := range nicknames() {
		if conn, err := kahoot.NewConn(gamePin); err != nil {
			fmt.Fprintln(os.Stderr, "failed to connect:", err)
			os.Exit(1)
		} else {
			defer conn.GracefulClose()
			conn.Login(nickname)
		}
	}

	fmt.Println("Kill this process to deauthenticate.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

func nicknames() []string {
	if len(os.Args) == 4 {
		count, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid count:", os.Args[3])
			os.Exit(1)
		}
		base := os.Args[2]
		res := make([]string, count)
		for x := 0; x < count; x++ {
			res[x] = base + strconv.Itoa(x+1)
		}
		return res
	}

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
	}

	return res
}
