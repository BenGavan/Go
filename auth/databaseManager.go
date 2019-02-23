package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	connStr = "postgres://postgres:password@localhost/sample_db?sslmode=disable"
)

func openDatabaseCon() (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		panic(err)
		return nil, err
	}
	fmt.Println("Database ping successful")
	return db, err
}

func isEmailValid(email string, c chan bool) {
	db, err := openDatabaseCon()
	rows, err := db.Query("SELECT EXISTS (SELECT TRUE FROM users WHERE email=$1);", email)
	if err != nil {
		fmt.Println("Error carrying out query")
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}()

	var result bool
	for rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("email exists:", result)
	}

	c <- !result
}

func isUsernameValid(username string, c chan bool) {
	db, err := openDatabaseCon()
	rows, err := db.Query("SELECT EXISTS (SELECT TRUE FROM users WHERE username=$1);", username)
	if err != nil {
		fmt.Println("Error carrying out query")
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}()

	var result bool
	for rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("username exists:", result)
	}

	c <- !result
}

func isPhoneNumberValid(number string, c chan bool) {
	db, err := openDatabaseCon()
	rows, err := db.Query("SELECT EXISTS (SELECT TRUE FROM users WHERE phone_number=$1);", number)
	if err != nil {
		fmt.Println("Error carrying out query")
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}()

	var result bool
	for rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println("phone number exists:", result)
	}

	c <- !result
}

func checkUid(uid string, c chan bool) {

	c <- true
}

func (user User) save() {
	db, err := openDatabaseCon()
	if err != nil {
		fmt.Println()
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()
	err = db.QueryRow(`INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		user.UID, user.DisplayName, user.Username, user.Email, user.PhoneNumber, user.Password, user.Registered).Scan()
	if err != nil {
		log.Println(err)
	}

	db.Close()
}

func (s session) save() {
	db, err := openDatabaseCon()
	if err != nil {
		fmt.Println()
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	err = db.QueryRow(`INSERT INTO sessions VALUES ($1, $2, $3, $4)`,
		s.uid, s.token, s.created, s.lastUsed).Scan()
	if err != nil {
		log.Println(err)
	}

	db.Close()
}

func getUIDFromSession(r *http.Request) string {
	token := getSessionTokenFromRequest(r)

	db, err := openDatabaseCon()
	rows, err := db.Query("SELECT uid FROM sessions WHERE token=$1;", token)
	if err != nil {
		fmt.Println("Error carrying out query")
		return ""
	}
	defer func() {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}()

	var uid string
	for rows.Next() {
		err := rows.Scan(&uid)
		if err != nil {
			log.Fatal(err)
			return ""
		}
		fmt.Println("uid from session token:", uid)
		return uid
	}
	return ""
}

func isTokenValid(token string) bool {
	db, err := openDatabaseCon()
	rows, err := db.Query("SELECT EXISTS (SELECT TRUE FROM sessions WHERE token=$1);", token)
	if err != nil {
		fmt.Println("Error carrying out query")
		return false
	}
	defer func() {
		if err := rows.Close(); err != nil {
			panic(err)
		}
	}()

	var result bool
	for rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			log.Fatal(err)
			return false
		}
		fmt.Println("phone number exists:", result)
	}
	return result
}
