package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

func main() {
	fmt.Printf("xml testing\n")

	err := unmarshal()
	if err != nil {
		panic(err)
	}
}

type Data struct {
	Outer outer `xml:"outer"`
}

type outer struct {
	Environment environment `xml:"environment"`
}

type environment struct {
	Temperature temperature `xml:"temperature"`
	Test        string      `xml:"test,attr"`
}

type temperature struct {
	Temperature string `xml:",chardata"`
	Type        string `xml:"type,attr"`
	Units       string `xml:"units,attr"`
}

func marshal() error {
	data := environment{Temperature: temperature{
		Temperature: "-12.3",
		Type:        "float",
		Units:       "c",
	}}

	xmlBytes, err := xml.Marshal(&data)
	if err != nil {
		return err
	}

	fmt.Printf("XML: %v\n", string(xmlBytes))

	return nil
}

func unmarshal() error {
	//s := "<environment><temperature units=\"c\">-12.3</temperature></environment>\n"
	//s := "<outer><environment><temperature units=\"c\">-12.3</temperature></environment></outer>\n"  // missing element attribute
	s := "<outer><environment test=\"test value\"><temperature units=\"c\">-12.3</temperature></environment></outer>\n" // full
	//s := "<outer><environment test=\"test value\"></environment></outer>\n" // Missing sub/child element

	sBytes := []byte(s)

	//top := Data{}
	var top outer
	err := xml.Unmarshal(sBytes, &top)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", structToString(top))

	return nil
}

func structToString(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Printf("Failed to convert struct to json")
		return ""
	}
	s := string(b)
	return s
}
