package kahoot

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

func gameSessionToken(gamePin int) (string, error) {
	resp, err := http.Get("https://kahoot.it/reserve/session/" + strconv.Itoa(gamePin))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	token := resp.Header.Get("X-Kahoot-Session-Token")
	var bodyObj struct {
		Challenge string `json:"challenge"`
	}
	if err := json.Unmarshal(body, &bodyObj); err != nil {
		return "", fmt.Errorf("parse session challenge: %s", err)
	}

	return decipherToken(token, bodyObj.Challenge)
}

func decipherToken(xToken, challenge string) (string, error) {
	r := bytes.NewReader([]byte(xToken))
	base64Dec := base64.NewDecoder(base64.StdEncoding, r)
	rawToken, err := ioutil.ReadAll(base64Dec)
	if err != nil {
		return "", fmt.Errorf("parse session token: %s", err)
	}

	challengeExpr := regexp.MustCompile("^\\(([0-9]*) \\+ ([0-9]*)\\) \\* ([0-9]*)$")
	match := challengeExpr.FindStringSubmatch(challenge)
	if match == nil {
		return "", fmt.Errorf("unsupported challenge: %s", challenge)
	}

	num1, _ := strconv.Atoi(match[1])
	num2, _ := strconv.Atoi(match[2])
	num3, _ := strconv.Atoi(match[3])
	mask := []byte(strconv.Itoa((num1 + num2) * num3))

	for i := range rawToken {
		rawToken[i] ^= mask[i%len(mask)]
	}

	return string(rawToken), nil
}
