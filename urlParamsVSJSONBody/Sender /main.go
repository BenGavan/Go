package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

/*
Makes the calls to the Receiver service (localhost:8010)
 */

type Request struct {
	ID string `json:"id"`
}

type Response struct {
	ID        string `json:"id"`
	TimeStamp int64  `json:"time_stamp"`
}

func main() {


	//timeParamCall()
	//timeJsonCall()

	//timeFuncCall(timeJsonCall, "json call")
	timeFuncCall(timeParamCall, "params call")
}

func timeParamCall() {
	startTime := time.Now()

	makeParamCall()

	endTime := time.Now()

	deltaTime := endTime.Sub(startTime)

	fmt.Printf("params call delta: %v\n", deltaTime.Nanoseconds())
}

func timeJsonCall() {
	startTime := time.Now()

	makeJSONCall()

	endTime := time.Now()

	deltaTime := endTime.Sub(startTime)

	fmt.Printf("json call delta: %v\n", deltaTime.Nanoseconds())
}

func timeFuncCall(f func(), description string) {
	startTime := time.Now()

	f()

	endTime := time.Now()

	deltaTime := endTime.Sub(startTime)

	fmt.Printf("%v delta: %v\n", description, deltaTime.Nanoseconds())
}

func makeParamCall() {
	id := randString(20)
	url := fmt.Sprintf( "http://localhost:8010/params?id=%v", id)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	handleResponse(resp, id)
}

func makeJSONCall() {
	id := randString(20)
	request := Request{ID: id}

	requestBodyBytes, err := json.Marshal(&request)
	if err != nil {
		panic(err)
	}

	requestBodyBytesBuffer := bytes.NewBuffer(requestBodyBytes)
	url :=  "http://localhost:8010/json"

	resp, err := http.Post(url, "application/json", requestBodyBytesBuffer)
	if err != nil {
		panic(err)
	}

	handleResponse(resp, id)
}

func handleResponse(resp *http.Response, id string) {
	defer resp.Body.Close()

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response Response
	err = json.Unmarshal(responseBodyBytes, &response)
	if err != nil {
		panic(err)
	}

	if response.ID != id {
		panic(errors.New("ids do not match"))
	}
}



//// Utils ////
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randString(length int) string {
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
