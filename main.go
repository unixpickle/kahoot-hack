package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
        "strings"
	"github.com/unixpickle/kahoot-hack/kahoot"
)

const ConcurrencyCount = 4

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
	if len(os.Args) == 4 {
		count, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid count:", os.Args[3])
			os.Exit(1)
		}
		base := os.Args[2]
		res := unicorn(15-len(base), count, base)
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
func unicorn(length, count int, base string) []string {
   char := []string{" ", " ", "᠎", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", "​", " ", " "}
   //charlist
   list := make([]string, 0)
   z := []string{"", "", ""}
   for i := 0; i < len(char); i++ {
      z[2] = char[i]
      z[0] = base
      list = append(list, strings.Join(z, ""))
   }
   for a := 0; a < length; a++ {
      for i := 0; i < len(list); i++ {
         for j := 0; j < len(char); j++ {
            z[2] = char[j]
            z[1] = list[i][len(z[0]):]
            z[0] = base
            if len(list) > count {
               break
            }
            list = append(list, strings.Join(z, ""))
         }
         if len(list) > count {
           break
         }
      }
      if len(list) > count {
         break
      }
   }
   return list
}
