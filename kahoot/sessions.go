package kahoot

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

var bruteForceErr = errors.New("not exactly one possible mask")

const tokenAttempts = 40

func gameSessionToken(gamePin int) (string, error) {
	for i := 0; i < tokenAttempts; i++ {
		token, err := attemptGameSessionToken(gamePin)
		if err != bruteForceErr {
			return token, err
		}
	}
	return "", errors.New("could not defeat session challenge")
}

func attemptGameSessionToken(gamePin int) (string, error) {
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

	maskNum, err := computeChallenge(challenge)
	if err != nil {
		maskNum, err = bruteForceChallenge(rawToken)
		if err != nil {
			return "", err
		}
	}
	mask := []byte(strconv.Itoa(maskNum))

	for i := range rawToken {
		rawToken[i] ^= mask[i%len(mask)]
	}

	return string(rawToken), nil
}

func computeChallenge(ch string) (int, error) {
	challengeExpr := regexp.MustCompile("^\\(([0-9]*)\\s*\\+\\s*([0-9]*)\\)\\s*\\*\\s*([0-9]*)$")
	match := challengeExpr.FindStringSubmatch(ch)
	if match != nil {
		num1, _ := strconv.Atoi(match[1])
		num2, _ := strconv.Atoi(match[2])
		num3, _ := strconv.Atoi(match[3])
		return (num1 + num2) * num3, nil
	}
	challengeExpr = regexp.MustCompile("^([0-9]*)\\s*\\*\\s*\\(([0-9]*)\\s*\\+\\s*([0-9]*)\\)$")
	match = challengeExpr.FindStringSubmatch(ch)
	if match != nil {
		num1, _ := strconv.Atoi(match[1])
		num2, _ := strconv.Atoi(match[2])
		num3, _ := strconv.Atoi(match[3])
		return num1 * (num2 + num3), nil
	}
	return 0, fmt.Errorf("unsupported challenge: %s", ch)
}

func bruteForceChallenge(rawToken []byte) (int, error) {
	var possibilities [][]byte
LengthLoop:
	for n := 1; n < 9; n++ {
		possible := make([]byte, n)
		for i := range possible {
			possible[i] = possibleMaskByte(rawToken, n, i)
			if possible[i] == 0 {
				continue LengthLoop
			}
		}
		possibilities = append(possibilities, possible)
	}
	if len(possibilities) != 1 {
		return 0, bruteForceErr
	}
	return strconv.Atoi(string(possibilities[0]))
}

func possibleMaskByte(rawToken []byte, chLen, byteIdx int) byte {
	start := 0
	if byteIdx == 0 {
		start = 1
	}
	possibs := []byte{}
PossibilityLoop:
	for b := start; b <= 9; b++ {
		numChar := byte(b) + '0'
		for i := byteIdx; i < len(rawToken); i += chLen {
			masked := rawToken[i] ^ numChar
			if !((masked >= 'a' && masked <= 'f') || (masked >= '0' && masked <= '9')) {
				continue PossibilityLoop
			}
		}
		possibs = append(possibs, numChar)
	}
	if len(possibs) != 1 {
		return 0
	}
	return possibs[0]
}
