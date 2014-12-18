package main

import (
	"fmt"
	"github.com/unixpickle/kahoot-hack/kahoot"
	"os"
)

func main() {
	var pin string
	var nickname string
	fmt.Print("Enter game pin: ")
	fmt.Scanln(&pin)
	fmt.Print("Enter nickname: ")
	fmt.Scanln(&nickname)
	fmt.Println("Connecting...")
	conn, err := kahoot.NewConnection(pin)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if err := conn.Handshake(); err != nil {
		fmt.Println("Handshake error:", err)
		os.Exit(1)
	}
	services := []string{"/service/player", "/service/controller",
		"/service/status"}
	err = conn.Subscribe(services)
	if err != nil {
		fmt.Println("Subscribe failed:", err)
		os.Exit(1)
	}
}
