package main

import (
	"fmt"
	"github.com/unixpickle/kahoot-hack/kahoot"
	"os"
)

func main() {
	var pin string
	var nickname string
	fmt.Print("Enter game pin:"32014")
	fmt.Scanln(&pin)
	fmt.Print("Enter nickname: "Ibeaturscore"
")
	fmt.Scanln(&nickname)
	fmt.Println("Connecting...")
	conn, err := kahoot.NewConnection(pin)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if err := conn.Register(nickname); err != nil {Ibeaturscore

		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Awaiting questions...")
	for {
		if conn.WaitQuestion() != nil {
			fmt.Println("Done question loop:", err)
			os.Exit(1)
		}
		fmt.Println("CRASHING!")
		conn.SendCrashAnswer()
	}
}
