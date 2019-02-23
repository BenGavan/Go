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
	"time"
)

type User struct {
	IsLoggedin bool   `json:"isLoggedIn"`
	Name       string `json:"name"`
}

func main() {
	timeTest()

	c := make(chan bool)
	cValue := make(chan int)

	go asyncCheck("test", c)
	go asynValue(cValue)

	asyncCheckResult := <- c
	asyncValue := <- cValue

	fmt.Println("AsyncCheckResult:", asyncCheckResult)
	fmt.Println("asyncVaklue:", asyncValue)

	//deleteFile("test.txt")
	//testFilePath := "testing.txt"
	//writeToFile(testFilePath, []byte("This now works"))
	//readFile(testFilePath)
	//createNode()
	http.Handle("/static", http.FileServer(http.Dir("/static/")))
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/hi", sayHello)
	http.HandleFunc("/", handleBase)
	http.HandleFunc("/servehtml", handleServeHTMLFile)
	http.HandleFunc("/serveimage", handleServeImage)
	http.HandleFunc("/user", handleServeJSONfromStruct)
	http.HandleFunc("/users", handleServeJSONArrayfromStruct)
	http.HandleFunc("/createcookie", handleCreateCookie)
	http.HandleFunc("/deletecookie", handleDeleteCookie)
	http.HandleFunc("/readcookie", handleGetCookieValue)
	http.HandleFunc("/basicform", handleBasicWebFormParse)
	http.HandleFunc("/upload", ReceiveFile)
	http.HandleFunc("/upload2", ReceiveMultiValueForm)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleBase(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fmt.Fprintln(w, "404: Page not found")
		return
	}
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "Hi :)")
		fmt.Println("GET")
	}
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
}

/*
////////////////   ------ chan ------   ////////////////
 */
 func asyncCheck(testString string, out chan bool) {
 	out <- true
 }

 func asynValue(out chan int) {
 	out <- 123
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
	defer file.Close()
	if err != nil {
		fmt.Println("file has not been created")
		return
	}
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
	text  := r.Form["x"]
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
