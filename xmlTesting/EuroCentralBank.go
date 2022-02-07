package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

type Envelope struct {
	XMLName xml.Name `xml:"gesmes:Envelope"`
	Subject string `xml:"gesmes:subject"`
}

func main() {
	fmt.Printf("European Central bank\n")

	err := marshal()
	if err != nil {
		panic(err)
	}
}

func marshal() error{

	data := Envelope{Subject: "Reference rates"}

	xmlBytes, err := xml.Marshal(&data)
	if err != nil {
		return err
	}

	xmlStringBytes, err := formatXML(xmlBytes)
	if err != nil {
		return err
	}

	fmt.Printf("Xml = %v\n", string(xmlStringBytes))

	return err
}


func formatXML(data []byte) ([]byte, error) {
	b := &bytes.Buffer{}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	encoder := xml.NewEncoder(b)
	encoder.Indent("", "  ")
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			encoder.Flush()
			return b.Bytes(), nil
		}
		if err != nil {
			return nil, err
		}
		err = encoder.EncodeToken(token)
		if err != nil {
			return nil, err
		}
	}
}