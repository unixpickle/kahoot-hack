package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/unixpickle/kahoot-hack/kahoot"
)

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Replace(r.URL.Path, "/", "", -1)
		if path == "" {
			http.ServeFile(w, r, "assets/index.html")
		} else if path == "flood" {
			handleFlood(w, r)
		} else {
			http.ServeFile(w, r, "assets/"+path)
		}
	})

	http.ListenAndServe(":"+os.Args[1], nil)
}

func handleFlood(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		w.Write([]byte("<!doctype html><head></head><body>Invalid form</body><html>"))
		return
	}
	pin := strings.TrimSpace(r.PostFormValue("pin"))
	nickname := r.PostFormValue("nickname")
	gamePin, _ := strconv.Atoi(pin)
	for i := 0; i < 20; i++ {
		conn, err := kahoot.NewConn(gamePin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to connect:", err)
			os.Exit(1)
		}
		conn.Login(nickname + strconv.Itoa(i+1))
		defer conn.Close()
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
