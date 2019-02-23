package main

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func main() {
	var (
		driver neo4j.Driver
		session neo4j.Session
		result neo4j.Result
		err error
	)

	if driver, err = neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("username", "password", "")); err != nil {
		panic(err) // handle error
	}
	// handle driver lifetime based on your application lifetime requirements
	// driver's lifetime is usually bound by the application lifetime, which usually implies one driver instance per application
	defer driver.Close()

	if session, err = driver.Session(neo4j.AccessModeWrite); err != nil {
		panic(err) // handle error
	}
	defer session.Close()

	result, err = session.Run("CREATE (n:Item { id: $id, name: $name }) RETURN n.id, n.name", map[string]interface{}{
		"id": 1,
		"name": "Item 1",
	})
	if err != nil {
		panic(err) // handle error
	}

	for result.Next() {
		fmt.Printf("Created Item with Id = '%d' and Name = '%s'\n", result.Record().GetByIndex(0).(int64), result.Record().GetByIndex(1).(string))
	}
	if err = result.Err(); err != nil {
		panic(err) // handle error
	}
}

func newDriver() (neo4j.Driver, error) {
	return neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("username", "password", ""))

	//driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("username", "password", "")); err != nil {
	//	return driver, err // handle error
	//}
	//return
}
