package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"math"
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
		res := make([]string, count)
		w := ""
		k := 0
		r := 0
		hi :=[20]string{" "," "," ","᠎"," "," "," "," "," "," "," "," "," "," "," "," ","　"," "," "," "}
		for x := 0; x < count; x++ {
			k =x
			w=""
			for y := int(math.Logb(float64(x))/4); y > 0; y-- {
					r = int(k)/int(math.Pow(20,float64(y)))
					k = int(k)-int(r)*int(math.Pow(20,float64(y)))
					w = hi[int(r)] + w
				}
			res[x] = base + w
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
