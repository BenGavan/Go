package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func main() {
	fmt.Printf("Hey\n")
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	//uri := "bolt://localhost:7687"
	//username := "neo4j"
	//password := "password"
	//greating, err := helloWorld(uri, username, password, true)
	//if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//}
	//fmt.Printf("Greeting: %v\n", greating)
}

type GraphDatabase struct {
	username string
	password string
	url      string
	driver   neo4j.Driver
}

func newGraphDatabase() (GraphDatabase, error) {
	g := GraphDatabase{
		username: "neo4j",
		password: "password",
		url:      "bolt://docker.for.mac.localhost:7687",
	}

	driver, err := g.newDriver()
	if err != nil {
		return g, err
	}
	g.driver = driver

	return g, err
}

func (g *GraphDatabase) newDriver() (neo4j.Driver, error) {
	config := func(conf *neo4j.Config) { conf.Encrypted = false }
	driver, err := neo4j.NewDriver(g.url, neo4j.BasicAuth(g.username, g.password, ""), config)
	return driver, err
}

func run() error {
	gdb, err := newGraphDatabase()
	if err != nil {
		return err
	}
	defer gdb.driver.Close()

	err = gdb.getItems()
	if err != nil {
		return err
	}

	return err
}

func (g *GraphDatabase) addNewItem() error {
	session, err := g.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	queryString := "CREATE (n:Item { id: $id, name: $name }) RETURN n.id as id, n.name as name"

	queryParams := map[string]interface{}{
		"id":   2,
		"name": "Item 2",
	}

	result, err := session.Run(queryString, queryParams)
	if err != nil {
		return err
	}

	for result.Next() {
		keys, err := result.Keys()
		fmt.Printf("Keys: %v\n", keys)
		if err != nil {
			return err
		}

		record := result.Record()
		fmt.Printf("Record: %v\n", record)

		resultSummary, err := result.Summary()
		if err != nil {
			return err
		}
		fmt.Printf("Result Summary: %v\n", resultSummary)

		id, ok := result.Record().Get("id")
		name, ok := result.Record().Get("name")
		if !ok {
			return errors.New("failed to extract result vales from result record")
		}

		fmt.Printf("ID = %v, name = %v\n", id, name)

		fmt.Printf("Created Item with Id = '%d' and Name = '%s'\n", result.Record().GetByIndex(0).(int64), result.Record().GetByIndex(1).(string))
	}
	return result.Err()
}

func (g *GraphDatabase) testCustomReturn() error {
	session, err := g.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	queryString := "CREATE (n:Item { id: $id, name: $name, age: $age }) RETURN {id: n.id, person: {name: n.name, age: n.age}}"

	queryParams := map[string]interface{}{
		"id":   "isString12334",
		"name": "Item 2",
		"age":  11,
	}

	result, err := session.Run(queryString, queryParams)
	if err != nil {
		return err
	}

	for result.Next() {
		keys, err := result.Keys()
		fmt.Printf("Keys: %v\n", keys)
		if err != nil {
			return err
		}

		recordValues := result.Record().Values()
		fmt.Printf("Record values: %v\n", recordValues)

		v := recordValues[0].(map[string]interface{})
		fmt.Printf("v: %v\n", v)

		idFromMap := v["id"]
		fmt.Printf("id from map: %v\n", idFromMap)

		personFromMap := v["person"].(map[string]interface{})
		fmt.Printf("Person from map: %v\n", personFromMap)

		nameFromPersonMap := personFromMap["name"]
		fmt.Printf("nameFromPersonMap: %v\n", nameFromPersonMap)

		resulrKeys := result.Record().Keys()
		fmt.Printf("resulrKeys: %v\n", resulrKeys)

		jsonStringBytes, err := json.Marshal(v)
		if err != nil {
			fmt.Printf("error marshling json from marshaling json to string: %v\n", err)
		}
		fmt.Printf("Json string: %v\n", string(jsonStringBytes))

		type Person struct {
			Name string `json:"name"`
			Age  int64  `json:"age"`
		}

		type Result struct {
			Id     string `json:"id"`
			Person Person `json:"person"`
		}

		var r Result
		err = json.Unmarshal(jsonStringBytes, &r)
		if err != nil {
			fmt.Printf("error unmarsheling json: %v\n", err)
		}
		fmt.Printf("Result from json unmarshling: ")
		PrintStruct(r)

		record := result.Record()
		fmt.Printf("Record: %v\n", record)

		resultSummary, err := result.Summary()
		if err != nil {
			return err
		}
		fmt.Printf("Result Summary: %v\n", resultSummary)

		id, ok := result.Record().Get("id")
		name, ok := result.Record().Get("name")
		if !ok {
			return errors.New("failed to extract result vales from result record")
		}

		fmt.Printf("ID = %v, name = %v\n", id, name)

		fmt.Printf("Created Item with Id = '%d' and Name = '%s'\n", result.Record().GetByIndex(0).(int64), result.Record().GetByIndex(1).(string))
	}
	return result.Err()
}

func (g *GraphDatabase) getItems() error {
	session, err := g.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	queryString := "MATCH (a:Item {id: 2}) RETURN {id: a.id, name: a.name}"

	queryParams := map[string]interface{}{

	}

	result, err := session.Run(queryString, queryParams)
	if err != nil {
		return err
	}

	////////////

	type Item struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}

	var items []Item

	for result.Next() {
		recordValues := result.Record().Values()
		fmt.Printf("Record values: %v\n", recordValues)

		v := recordValues[0].(map[string]interface{})
		fmt.Printf("v: %v\n", v)

		resulrKeys := result.Record().Keys()
		fmt.Printf("resulrKeys: %v\n", resulrKeys)

		jsonStringBytes, err := json.Marshal(v)
		if err != nil {
			fmt.Printf("error marshling json from marshaling json to string: %v\n", err)
		}
		fmt.Printf("Json string: %v\n", string(jsonStringBytes))

		var item Item
		err = json.Unmarshal(jsonStringBytes, &item)
		if err != nil {
			fmt.Printf("error unmarsheling json: %v\n", err)
		}
		items = append(items, item)

		fmt.Printf("Result from json unmarshling: ")
		PrintStruct(item)
	}

	PrintStruct(items)
	return err
}

func min() error {
	// configForNeo4j35 := func(conf *neo4j.Config) {}
	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }

	driver, err := neo4j.NewDriver("bolt://docker.for.mac.localhost:7687", neo4j.BasicAuth("neo4j", "password", ""), configForNeo4j40)
	PrintStruct(driver)
	if err != nil {
		return err
	}
	// handle driver lifetime based on your application lifetime requirements
	// driver's lifetime is usually bound by the application lifetime, which usually implies one driver instance per application
	defer driver.Close()

	// For multidatabase support, set sessionConfig.DatabaseName to requested database

	session, err := driver.Session(neo4j.AccessModeWrite)
	fmt.Printf("Session: %v, error: %v\n", session, err)
	PrintStruct(session)
	if err != nil {
		return err
	}
	defer session.Close()

	result, err := session.Run("CREATE (n:Item { id: $id, name: $name }) RETURN n.id, n.name", map[string]interface{}{
		"id":   1,
		"name": "Item 1",
	})
	fmt.Printf("Result: %v, error: %v\n", result, err)
	PrintStruct(result)
	if err != nil {
		return err
	}

	for result.Next() {
		fmt.Printf("Created Item with Id = '%d' and Name = '%s'\n", result.Record().GetByIndex(0).(int64), result.Record().GetByIndex(1).(string))
	}
	return result.Err()
}

func helloWorld(uri, username, password string, encrypted bool) (string, error) {
	authToken := neo4j.BasicAuth(username, password, "")

	fmt.Printf("Auth Token: %v\n", authToken)
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = encrypted
	})
	if driver == nil {
		fmt.Printf("driver is also nil")
	}
	if err != nil {
		return "", err
	}
	defer driver.Close()

	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return "", err
	}
	defer session.Close()

	greeting, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
			map[string]interface{}{"message": "hello, world"})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return "", err
	}

	return greeting.(string), nil
}

func PrintStruct(data interface{}) {
	s := StructToString(data)
	fmt.Println(s)
}

func StructToString(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Printf("Failed to convert struct to json")
		return ""
	}
	s := string(b)
	return s
}
