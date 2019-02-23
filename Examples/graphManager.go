package main

import (
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

/*
////////////////  ----- Neo4j  ---- ////////////////
 */
const (
	URI          = "bolt://neo4j:1234@localhost:7687"
	CreateNode   = "CREATE (n:NODE {foo: {foo}, bar: {bar}})"
	GetNode      = "MATCH (n:NODE) RETURN n.foo, n.bar"
	RelationNode = "MATCH path=(n:NODE)-[:REL]->(m) RETURN path"
	DeleteNodes  = "MATCH (n) DETACH DELETE n"
)

func createConnection() bolt.Conn {
	driver := bolt.NewDriver()
	con, err := driver.OpenNeo(URI)
	handleError(err)
	return con
}

func createNode() {
	con := createConnection()
	defer con.Close()

	st := prepareStatement(CreateNode, con)
	executeStatement(st)
}

// Prepare a new statement. This gives us the flexibility to
// cancel that statement without any request sent to Neo
func prepareStatement(query string, con bolt.Conn) bolt.Stmt {
	st, err := con.PrepareNeo(query)
	handleError(err)
	return st
}

// Executing a statement just returns summary information
func executeStatement(st bolt.Stmt) {
	result, err := st.ExecNeo(map[string]interface{}{"uid": "1hvgdN8e4", "name": "Ben"})
	handleError(err)
	numResult, err := result.RowsAffected()
	handleError(err)
	fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 1

	// Closing the statement will also close the rows
	st.Close()
}



// Simple function to  take care of errors,
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}



