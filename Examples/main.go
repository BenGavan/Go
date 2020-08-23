package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

type User struct {
	IsLoggedin bool   `json:"isLoggedIn"`
	Name       string `json:"name"`
}

func main() {
	myDateString := "2000-01-08 16:30"
	fmt.Println("My Starting Date:\t", myDateString)

	// Parse the date string into Go's time object
	// The 1st param specifies the format, 2nd is our date string
	myDate, err := time.Parse("2006-01-02 15:04", myDateString)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Data: %v\n", myDate)
	unix := myDate.UTC().Unix()
	fmt.Printf("Unix: %v\n", unix)

	fmt.Printf("Time now: %v\n", time.Now().UTC().Unix())
	fmt.Printf("Time now: %v\n", time.Now().UTC().UnixNano())
}

func writeColumnsToFileTest() {
	a1 := make([]float64, 3)
	a1[0] = 11
	a1[1] = 12
	a1[2] = 13

	a2 := make([]float64, 3)
	a2[0] = 21
	a2[1] = 22
	a2[2] = 23

	a3 := make([]float64, 3)
	a3[0] = 31
	a3[1] = 32
	a3[2] = 33

	path := "testcol.txt"

	writeColumnsToFile(path, a1, a2, a3)
}

func run() {
	timeTest()
	//deleteFile("test.txt")
	//testFilePath := "testing.txt"
	//writeToFile(testFilePath, []byte("This now works"))
	//readFile(testFilePath)
	//createNode()
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static/"))
	mux.Handle("/static", http.FileServer(http.Dir("/static/")))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.Handle("/favicon.ico", http.NotFoundHandler())

	mux.HandleFunc("/hi", sayHello)
	mux.HandleFunc("/", handleBase)
	mux.HandleFunc("/servehtml", handleServeHTMLFile)
	mux.HandleFunc("/serveimage", handleServeImage)
	mux.HandleFunc("/unsafetemplate", handleUnsafeTemplate)
	mux.HandleFunc("/safetemplate", handleSafeTemplate)
	mux.HandleFunc("/user", handleServeJSONfromStruct)
	mux.HandleFunc("/users", handleServeJSONArrayfromStruct)
	mux.HandleFunc("/createcookie", handleCreateCookie)
	mux.HandleFunc("/deletecookie", handleDeleteCookie)
	mux.HandleFunc("/readcookie", handleGetCookieValue)
	mux.HandleFunc("/basicform", handleBasicWebFormParse)
	mux.HandleFunc("/upload", ReceiveFile)
	mux.HandleFunc("/upload2", ReceiveMultiValueForm)

	mux.HandleFunc("/testasnc", handleTestAsync)

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

func multipleParameters(arrs ...[]float64) {
	fmt.Println(arrs)
	for i := 0; i < len(arrs); i++ {
		fmt.Println(arrs[i])
	}
}

//func registerRoutes(mux *http.ServeMux) {
//	mux.Handle()
//}

func handleBase(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	if r.URL.Path != "/" {
		fmt.Fprintln(w, "404: Page not found")
		return
	}
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "Hi :)")
	}
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	_, err := w.Write([]byte(message))
	if err != nil {
		handleError(err)
	}
}

/*
append
 */
func appendExamples() {
	var ts []byte
	var exs []byte
	ts = append(ts, exs...)
	ts = append(ts, 1, 3, 3)
	ts = append(ts, "string"...)
}

/*
////////////////   ------ chan / async ------   ////////////////
*/
func asyncCheck(testString string, out chan bool) {
	out <- true
}

func asyncValue(out chan int) {
	out <- 123
}

func asyncString(inputString string, out chan string) {
	out <- inputString
}

func handleTestAsync(w http.ResponseWriter, r *http.Request) {
	a, b, c := asyncTesting()

	aString := fmt.Sprint(a)
	bString := string(b)

	var outBytes [][]byte

	outBytes = append(outBytes, []byte(aString))
	outBytes = append(outBytes, []byte(bString))
	outBytes = append(outBytes, []byte(c))

	var flatOutBytes []byte

	for i := 0; i < len(outBytes); i++ {
		currentBytes := outBytes[i]
		for j := 0; j < len(currentBytes); i++ {
			flatOutBytes = append(flatOutBytes, currentBytes[j])
		}
	}

	_, err := w.Write(flatOutBytes)
	handleError(err)
}

func asyncTesting() (bool, int, string) {
	c := make(chan bool)
	cValue := make(chan int)
	cStringValue := make(chan string)

	go asyncCheck("test", c)
	go asyncValue(cValue)
	go asyncString("This is the hard coded parameter string.", cStringValue)

	asyncCheckResult := <-c
	asyncValue := <-cValue
	asyncStringValue := <-cStringValue

	fmt.Println("AsyncCheckResult:", asyncCheckResult)
	fmt.Println("asyncVaklue:", asyncValue)
	fmt.Println("asyncString value", asyncStringValue)

	close(c)
	close(cValue)
	close(cStringValue)

	//return <-c, <-cValue, <-cStringValue
	return asyncCheckResult, asyncValue, asyncStringValue
}

////////////////  ------ html/template ----- ////////////////
func handleUnsafeTemplate(w http.ResponseWriter, r *http.Request) {
	// The same as safe template; however, we are using temple from the Text package (text/template)
	t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	err = t.ExecuteTemplate(w, "T", "<script>alert('you have been pwned')</script>;")
	handleError(err)
}

func handleSafeTemplate(w http.ResponseWriter, r *http.Request) {
	// The same as unsafe template; however, we are using temple from the HTML package (html/template)
	t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	handleError(err)
	err = t.ExecuteTemplate(w, "T", "<script>alert('you have been pwned')</script>;")
	handleError(err)
}

/*
////////////////  ------  Files  ------  ////////////////
*/
///// Serve file
func handleServeHTMLFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// Serve Image
func handleServeImage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "images/test.jpeg")
}

// Create file/dirrectory if it doesn't exist
func createFile(path string) {
	if fileDoesExist(path) {
		fmt.Println("File", path, "already exists")
		return
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("file has not been created")
		return
	}
	defer file.Close()
	fmt.Println("file created")
}

func writeFile(path string) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Error(err)
		return
	}
	defer file.Close()

	_, err = file.Write([]byte("hi"))
	if err != nil {
		fmt.Println(err)
		return
	}

	err = file.Sync()
	if err != nil {
		panic(err)
		return
	}

	fmt.Println("File successfully uploaded")
}

func writeColumnsToFile(path string, cols ...[]float64) {
	if !isArraysOfEqualLength(cols) {
		fmt.Println("ERROR: Arrays not of equal length")
		return
	}
	var outString string
	for i := 0; i < len(cols[0]); i++ {
		for j := range cols {
			outString += fmt.Sprint(cols[j][i], "\t")
		}
		outString += "\n"
	}
	outBytes := []byte(outString)
	writeToFile(path, outBytes)
}

func writeToFile(path string, bytes []byte) {
	//fmt.Println("path:", path, "bytes:", bytes)
	createFile(path)
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Error(err)
		return
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = file.Sync()
	if err != nil {
		panic(err)
		return
	}

	fmt.Println("File successfully uploaded")
}

func readBytesFromFile(filepath string) ([]byte, error) {
	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer

	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

func readFile(path string) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var bytes []byte
	var text = make([]byte, 1)
	for {
		_, err = file.Read(text)

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}

		bytes = append(bytes, text[0])
	}

	//var bytes []byte
	//_, err = file.Read(bytes)
	//fmt.Println(bytes)
	//outString := string(bytes)
	textString := string(bytes)
	fmt.Println("Text from File:", textString)
}

// Delete file
func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("File:", path, "Successfully deleted")
}

// Check if file exists
func fileDoesExist(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

func fileDoesNotExist(path string) bool {
	return !fileDoesExist(path)
}

/*
////////////////  ----- JSON -----  ////////////////
*/

func handleServeJSONfromStruct(w http.ResponseWriter, r *http.Request) {
	user := User{true, "Ben Gavan"}
	json.NewEncoder(w).Encode(&user)
}

func handleServeJSONArrayfromStruct(w http.ResponseWriter, r *http.Request) {
	userOne := User{true, "Ben Gavan"}
	userTwo := User{false, "Some name"}
	userThree := User{true, "Matt Singleton"}

	var users = []User{userOne, userTwo}

	users = append(users, userThree)

	json.NewEncoder(w).Encode(&users)
}

/*
////////////////  ---- Cookies ----   //////////////////
*/

func handleGetCookieValue(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("cookiename")

	if err != nil {
		w.Write([]byte("error in reading cookie : " + err.Error() + "\n"))
		return
	}
	value := c.Value
	w.Write([]byte("cookie has : " + value + "\n"))
}

func handleCreateCookie(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   "cookiename",
		Value:  "thisisthecookie'svalue",
		MaxAge: 360000}
	http.SetCookie(w, &c)

	w.Write([]byte("new cookie created!\n"))
}

func handleDeleteCookie(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   "cookiename",
		MaxAge: -1}
	http.SetCookie(w, &c)

	w.Write([]byte("old cookie deleted!\n"))
}

/*
////////////////  -----  Parse Form Data  -----  ////////////////
*/
func handleBasicWebFormParse(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	switch r.Method {
	case http.MethodGet:
		//fmt.Println(time.Now(), "|", r.URL.Path, "|", r.Method, "|")
	case http.MethodPost:
		v := r.FormValue("q")
		fmt.Println("the value is: ", v)
	}
	v := r.FormValue("q")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `
	<form method="POST">
	 <input type="text" name="q">
	 <input type="submit">
	</form>
	<br>`+v)
}

func handleBasicIOSFormParse(w http.ResponseWriter, r *http.Request) {

}

func handleCustomPostForm(w http.ResponseWriter, r *http.Request) {

}

func handleImageUpload(w http.ResponseWriter, r *http.Request) {

}

func Cleaner(w http.ResponseWriter, r *http.Request) {
	// Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var u User
	err = json.Unmarshal(b, &u)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal into raw string-interface map
	var raw map[string]interface{}
	err = json.Unmarshal(b, &raw)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	t := raw["username"]
	fmt.Println("username received:", t)

	output, err := json.Marshal(u)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

//
func ReceiveFile(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	var Buf bytes.Buffer
	// in your case file would be fileupload
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fmt.Println(header)
	//name := strings.Split(header.Filename, ".")
	//fmt.Printf("File name %s\n", name[0])
	fmt.Printf("File name %s\n", "image.jpg")
	// Copy the file data to my buffer
	io.Copy(&Buf, file)
	// do something with the contents...
	// I normally have a struct defined and unmarshal into a struct, but this will
	// work as an example
	writeToFile("test.jpg", Buf.Bytes())

	contents := Buf.String()
	fmt.Println(contents)
	// I reset the buffer in case I want to use it again
	// reduces memory allocations in more intense projects
	Buf.Reset()
	// do something else
	// etc write header
	fmt.Fprint(w, "Success")
	return
}

func ReceiveMultiValueForm(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	var Buf bytes.Buffer
	r.ParseForm()
	//t := r.MultipartForm.Value
	//fmt.Println(t)
	// in your case file would be fileupload
	text := r.Form["x"]
	fmt.Println("The textis:", text)

	file, _, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//name := strings.Split(header.Filename, ".")
	//fmt.Printf("File name %s\n", name[0])
	fmt.Printf("File name %s\n", "image.jpg")
	// Copy the file data to my buffer
	io.Copy(&Buf, file)
	// do something with the contents...
	// I normally have a struct defined and unmarshal into a struct, but this will
	// work as an example
	writeToFile("test.jpg", Buf.Bytes())

	for key, value := range r.Form {
		fmt.Println(key, value)
	}

	//contents := Buf.String()
	//fmt.Println(contents)
	// I reset the buffer in case I want to use it again
	// reduces memory allocations in more intense projects
	Buf.Reset()
	// do something else
	// etc write header
	fmt.Fprint(w, "Success")
	return
}

/*
////////////////  ----- Time  ---- ////////////////
*/
func timeTest() {
	currentTime := time.Now()
	currentUNIX := time.Now().Unix()
	reformedTime := time.Unix(currentUNIX, 0)

	fmt.Println("Current Time:", currentTime)
	fmt.Println("Current UNIX time:", currentUNIX)
	fmt.Println("Reformed time from unix:", reformedTime)
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

func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
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

/*
////////////////  ----- Postgresql  ---- ////////////////
*/

/*
////////////////  ----- Util functions (aka random useful shit)  ---- ////////////////
*/
func printInfo(r *http.Request) {
	fmt.Println(time.Now(), "|", r.URL.Path, "|", r.Method, "|")
}

func printX() {
	for i := 0; i < 51; i++ {
		fmt.Print("X")
	}
}

//func handleError(err error) {
//	if err != nil {
//		log.Error(err)
//	}
//}

func isArraysOfEqualLength(arrays [][]float64) bool {
	for i := 0; i < len(arrays)-1; i++ {
		if len(arrays[i]) != len(arrays[i+1]) {
			return false
		}
	}
	return true
}
