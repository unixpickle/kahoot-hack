package kahoot

import (
	"net/http"
	"strconv"
)

func gameSessionToken(gamePin int) (string, error) {
	resp, err := http.Get("https://kahoot.it/reserve/session/" + strconv.Itoa(gamePin))
	if err != nil {
		return "", err
	}
	token := resp.Header.Get("X-Kahoot-Session-Token")
	resp.Body.Close()
	return token, nil
}
