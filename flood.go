package main

import (
	"fmt"
	"github.com/unixpickle/kahoot-hack/kahoot"
	"os"
	"strconv"
)

func main() {
	var pin string
	var nickname string
	var numNames string
	fmt.Print("Enter game pin: ")
	fmt.Scanln(&pin)
	fmt.Print("Enter nickname: ")
	fmt.Scanln(&nickname)
	fmt.Print("   Enter count: ")
	fmt.Scanln(&numNames)
	count, err := strconv.Atoi(numNames)
	if err != nil {
		fmt.Println("Invalid count:", numNames)
		os.Exit(1)
	}
	for i := 0; i < count; i++ {
		istr := strconv.Itoa(i + 1)
		newNick := nickname + istr
		fmt.Println("Registering", newNick, "...")
		conn, err := kahoot.NewConnection(pin)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		if err := conn.Register(newNick); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		go func() {
			for {
				_, err := conn.Read()
				if err != nil {
					return
				}
			}
		}()
	}
	fmt.Println("Done. Hit enter to exit.")
	fmt.Scanln()
}
