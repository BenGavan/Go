package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type User2 struct {
	uid int
	username string
}

func mainDatabase() {
	//connStr := "user=postgres dbname=sample_db sslmode=disable"
	connStr := "postgres://postgres:password@localhost/sample_db?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Now connected to database")

	//age := 21
	rows, err := db.Query("SELECT * FROM test2")
	if err != nil {
		fmt.Println("Error carrying out query")
		return
	}
	defer rows.Close()

	users := make([]User2, 0)
	for rows.Next() {
		user := User2{}
		err := rows.Scan(&user.uid, &user.username)
		if err != nil {
			log.Fatal(err)
			return
		}
		users = append(users, user)
	}

	fmt.Println(users)



	err = db.QueryRow(`INSERT INTO test2 VALUES (123,'ben123')`).Scan()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	db.Close()

	//insertRow()
}

func insertRow() {
	connStr := "user=localhost dbname=sample_db sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	//var userid int
	//err = db.QueryRow(`INSERT INTO users(name, favorite_fruit, age) VALUES('beatrice', 'starfruit', 93) RETURNING id`).Scan(&userid)
	err = db.QueryRow(`INSERT INTO users(name) VALUES('ben') `).Scan()
	if err != nil {
		log.Fatal(err)
	}
}
