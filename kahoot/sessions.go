package kahoot

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

const sessionTokenTimeout = time.Minute

var sessionTokenExpiration time.Time
var sessionTokenLock sync.RWMutex
var sessionTokenCache string
var sessionTokenCacheGamePin int = -1

func gameSessionToken(gamePin int) (string, error) {
	sessionTokenLock.RLock()
	if sessionTokenCacheGamePin != gamePin || time.Now().After(sessionTokenExpiration) {
		sessionTokenLock.RUnlock()
		return requestNewSessionToken(gamePin)
	}
	defer sessionTokenLock.RUnlock()
	return sessionTokenCache, nil
}

func requestNewSessionToken(gamePin int) (string, error) {
	sessionTokenLock.Lock()
	defer sessionTokenLock.Unlock()

	// May occur if two goroutines request a new session token simultaneously.
	if sessionTokenCacheGamePin == gamePin && !time.Now().After(sessionTokenExpiration) {
		return sessionTokenCache, nil
	}

	resp, err := http.Get("https://kahoot.it/reserve/session/" + strconv.Itoa(gamePin))
	if err != nil {
		return "", err
	}
	token := resp.Header.Get("X-Kahoot-Session-Token")
	resp.Body.Close()
	sessionTokenExpiration = time.Now().Add(sessionTokenTimeout)
	sessionTokenCacheGamePin = gamePin
	sessionTokenCache = token
	return token, nil
}
