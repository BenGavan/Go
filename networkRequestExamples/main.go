package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func customPostRequest() {
	data := map[string]interface{}{
		"name": "Ben Gavan",
		"email": "ben@gmail.com",
	}

	requestBody, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}

	requestBuffer := bytes.NewBuffer(requestBody)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	url := ""

	request, err := http.NewRequest(http.MethodPost, url, requestBuffer)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("Content-type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}

func fileUpload() {
	filepath := ""
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// New Buffer for request bytes
	var requestBody bytes.Buffer

	// To multi-part form writer (will write the necessary bytes to requestBodyBuffer)
	multipartWriter := multipart.NewWriter(&requestBody)

	// Add file to form
	fileWriter, err := multipartWriter.CreateFormFile("file_field", "filemane.txt")
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		log.Fatal(err)
	}

	// Add another normal field to form
	fieldWriter, err := multipartWriter.CreateFormField("normal_field")
	if err != nil {
		log.Fatal(err)
	}

	_, err = fieldWriter.Write([]byte("another value"))
	if err != nil {
		log.Fatal(err)
	}

	// Close the multipart file/field - also writes the ending boundary
	err = multipartWriter.Close()
	if err != nil {
		log.Fatal(err)
	}

	url := ""

	// Make POST request
	request, err := http.NewRequest(http.MethodPost, url, &requestBody)
	if err != nil {
		log.Fatal(err)
	}

	// Set Headers
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// Carry out request
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout, // In example, this is not included
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)

}

func run() error {
	return nil
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
