package kahoot

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

var (
	challengeRegexp = regexp.MustCompile(`^decode\.call\(this, '([a-zA-Z0-9]*)'\); ` +
		`function decode\(message\) \{var offset = ([0-9\+\*\(\)\s]*); ` +
		`if \(this\.angular\.[a-zA-Z]*\(offset\)\) \{` +
		`console.log\("Offset derived as: \{", offset, "\}"\);\}` +
		`return _\.replace\(message, /\./g, function\(char, position\) ` +
		`\{return String\.fromCharCode\(\(\(\(char\.charCodeAt\(0\) \* position\)` +
		` \+ offset\) % 77\) \+ 48\);\}\);\}$`)
)

func gameSessionToken(gamePin string) (string, error) {
	return attemptGameSessionToken(gamePin)
}

func attemptGameSessionToken(gamePin string) (string, error) {
	resp, err := http.Get("https://kahoot.it/reserve/session/" + gamePin)
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
		if string(body) == "Not found" {
			return "", fmt.Errorf("game pin not found: %s", gamePin)
		}
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

	mask, err := computeChallenge(challenge)
	if err != nil {
		return "", errors.New("failed to defeat challenge: " + challenge)
	}

	for i := range rawToken {
		rawToken[i] ^= mask[i%len(mask)]
	}

	return string(rawToken), nil
}

func computeChallenge(ch string) ([]byte, error) {
	submatch := challengeRegexp.FindStringSubmatch(ch)
	if submatch != nil {
		offset, err := eval(submatch[2])
		if err == nil {
			var newRunes []rune
			for i, x := range submatch[1] {
				n := (((int64(x) * int64(i)) + offset) % 77) + 48
				newRunes = append(newRunes, rune(n))
			}
			return []byte(string(newRunes)), nil
		}
	}

	evalURL := url.URL{
		Scheme:   "http",
		Host:     "safeval.pw",
		Path:     "/eval",
		RawQuery: url.Values{"code": []string{ch}}.Encode(),
	}
	resp, err := http.Get(evalURL.String())
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("server failed to evaluate: " + ch)
	}
	return ioutil.ReadAll(resp.Body)
}
