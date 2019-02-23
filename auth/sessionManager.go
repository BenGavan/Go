package main

import (
	"fmt"
	"net/http"
	"time"
)

type session struct {
	uid      string
	token    string
	created  int64
	lastUsed int64
}

const (
	SESSION_TOKEN = "session_token"
)

func (user User) createSession(w http.ResponseWriter) {

	token := createSessionToken()
	cookie := createSessionCookie(token)
	created := time.Now().Unix()

	session := session{user.UID, token, created, created}
	session.save()

	http.SetCookie(w, &cookie)

	fmt.Println("Session:", session)
}

func (user User) removeSession(w http.ResponseWriter) {
	// TODO remove session from session table

	c := http.Cookie{
		Name:   SESSION_TOKEN,
		MaxAge: -1}
	http.SetCookie(w, &c)
}

func (s session) checkIfValid() bool {
	return false
}

func isSessionValid(r *http.Request) bool {
	token := getSessionTokenFromRequest(r)
	return isTokenValid(token)
}

func createSessionCookie(token string) http.Cookie {
	return http.Cookie{
		Name:   SESSION_TOKEN,
		Value:  token,
		MaxAge: 360000}
}

func createSessionToken() string {
	return RandString(15)
}

func getSessionTokenFromRequest(r *http.Request) string {
	c, err := r.Cookie(SESSION_TOKEN)
	if err != nil {
		return ""
	}
	return c.Value
}
