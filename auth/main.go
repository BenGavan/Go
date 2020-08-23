package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type User struct {
	UID               string `json:"uid,omitempty"`
	Username          string `json:"username,omitempty"`
	Email             string `json:"email,omitempty"`
	PhoneNumber       string `json:"phone_number,omitempty"`
	DisplayName       string `json:"display_name,omitempty"`
	Password          string `json:"password,omitempty"`
	DateOfBirth       string `json:"date_of_birth,omitempty"`
	Registered        int64  `json:"registered,omitempty"`
	NumberOfFollowers int    `json:"number_of_followers,omitempty"`
	NumberFollowing   int    `json:"number_following,omitempty"`
}

type RegisterResponse struct {
	UID                    string `json:"uid,omitempty"`
	IsSuccessFul           bool   `json:"isUserRegistrationSuccessful,omitempty"`
	IsUsernameValid        bool   `json:"isUsernameValid,omitempty"`
	IsUsernameAvailable    bool   `json:"isUsernameAvailable,omitempty"`
	IsPhoneNumberValid     bool   `json:"isPhoneNumberValid,omitempty"`
	IsPhoneNumberAvailable bool   `json:"isPhoneNumberAvailable,omitempty"`
	IsEmailTakenValid      bool   `json:"isEmailValid,omitempty"`
	IsEmailTakenAvailable  bool   `json:"isEmailAvailable,omitempty"`
}

func main() {
	fmt.Println(":)")

	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/test", handleTestProtectedRoot)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	printInfo(r)

	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			fmt.Println("something went wrong parsing form")
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer r.Body.Close()

		// Unmarshal body data
		var user User
		err = json.Unmarshal(b, &user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		uidChan := make(chan string)
		defer close(uidChan)
		emailCheckChan := make(chan bool)
		defer close(emailCheckChan)
		usernameCheckChan := make(chan bool)
		defer close(usernameCheckChan)
		phoneNumberChan := make(chan bool)
		defer close(phoneNumberChan)

		go newUID(uidChan)
		go isEmailValid(user.Email, emailCheckChan)
		go isUsernameValid(user.Username, usernameCheckChan)
		go isPhoneNumberValid(user.PhoneNumber, phoneNumberChan)

		user.UID = <-uidChan
		isEmailValid := <-emailCheckChan
		isUsernameValid := <-usernameCheckChan
		isPhoneNumberValid := <-phoneNumberChan

		fmt.Println("User is:", user)

		fmt.Println("UID:", user.UID)
		fmt.Println("isEmailValid:", isEmailValid)
		fmt.Println("isUsernameValid:", isUsernameValid)
		fmt.Println("isPhoneNumberValid:", isPhoneNumberValid)

		if !isUsernameValid || !isPhoneNumberValid || !isEmailValid {
			resp := RegisterResponse{user.UID, false, isUsernameValid,
				isUsernameValid, isPhoneNumberValid, isPhoneNumberValid,
				isEmailValid, isEmailValid}
			err = json.NewEncoder(w).Encode(resp)
			if err != nil {
				log.Fatal(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println(resp)
			return
		}

		user.Registered = time.Now().Unix()
		time.Now().UnixNano()

		user.save()
		user.createSession(w)

		resp := RegisterResponse{user.UID, true, isUsernameValid,
			isUsernameValid, isPhoneNumberValid, isPhoneNumberValid,
			isEmailValid, isEmailValid}
		fmt.Println(user)
		fmt.Println(resp)
		fmt.Println(resp.IsSuccessFul)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			fmt.Println("something went wrong parsing form")
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer r.Body.Close()

		var user User
		err = json.Unmarshal(b, &user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
}

func handleTestProtectedRoot(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	if !isSessionValid(r) {
		http.Error(w, "User Unauthorised", http.StatusUnauthorized)
		print("User Unauthorised:", http.StatusUnauthorized)
		return
	}
	print("User is authorised")

	uid := getUIDFromSession(r)
	print("UID:", uid)
}

func newUID(c chan string) {
	// TODO ensure uid generated is unique
	c <- createUID()
}

func createUID() string {
	return RandString(15)
}

func hashPassword(password string) (string, string) {
	salt := RandString(8)
	hash := hashPasswordWithSalt(password, salt)
	return hash, salt
}

func hashPasswordWithSalt(password string, salt string) string {
	passWithSalt := salt + password
	h := sha256.New()
	h.Write([]byte(passWithSalt))
	hash := string(h.Sum(nil))
	return hash
}

func printInfo(r *http.Request) {
	fmt.Println(time.Now(), "|", r.URL.Path, "|", r.Method, "|")
}

func printSeparator() {
	fmt.Println(strings.Repeat("-", 20))
}

/*
////////////////  ----- Random String  ---- ////////////////
*/
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(length int) string {
	b := make([]byte, length)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := length-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
