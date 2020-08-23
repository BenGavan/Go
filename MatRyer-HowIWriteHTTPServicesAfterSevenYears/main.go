package main

import (
	"encoding/json"
	"fmt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/errors"
	"html/template"
	"net/http"
	"os"
	"sync"
)

/*
In Summary:
 - Write boring code
 - 'server' struct to hold dependencies (avoid global state/variables"
 - All routing in one place line 'route.go'
 - Handlers are anonymous funcs returned from methods
 - Make use of closured environment (request/response types, setup dependencies, lazy initialisation, etc)
 - use httptest for testing, anf just call the methods
 */

/*
Testing
 - net/http/httptest is really good - check it out - it's a standard package.
 */

/*
If the server type/object gets too big and noisy, we can have many servers.
 - Enables to separate out dependencies = very clear (a story telling point of view)
 - easy to see what's going on at a glance
Eg. a "people" server in people.go and a "comments" server in comments.go
 */
type server struct {
	db *someDatabase
	//router *somerouter
	router *http.ServeMux
	email EmailSender
}

func newServer() *server {
	// Only setup things that are the server, eg routes
	// Don't set up dependencies like database connections or sorting out the logger
	// We wnat to be able to use this constructor in different ways (so only set up the tings that are required for it to be a server)
	// Can pass in dependencies if there's like one but best not (perhaps pass the server to the initializer of those dependencies??)
	s := &server{}
	s.routes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Keep this like this... it would bee weird to do anything else
	// If we want to do something like logging, use middle ware
	s.router.ServeHTTP(w, r)
}

// Handlers hang off the server
// Allows them to access the dependencies
// Other handlers have access to s too, so have to  be careful with data races.
func (s *server) handleSomething() http.HandlerFunc {
	// Put some programming here
}
// Naming of handler methods:
/*
Since autocomplete loists and docs are alphabetically sorted -  group related to functionality
handlerTasksCreate
handlerTasksDone
handlerTasksGet

handleAuthLogin
handleAuthLogout
 */

/*
Return the handler
 - Allows for handler-specific setup on server startup
Take Arguments for handler-specific dependencies.
 - This  makes vary clear what the handler need to do its job
 - Type safety and compile-time checks
e.g:
 - *template.Template
 - *rand.Rand = (so in testing can pass a known random number generator (so known seed thus know output) so can be used in testing)
 */
func (s *server) handleSomethingElse(format string) http.HandlerFunc {
	thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request) {
		// use thing
		fmt.Fprintf(w, format, r.FormValue("name"))
	}
}

/*
Request and response data types
 - Story Telling here - they are only declared here so obviously only use/required here
 - colocate stuff is easier to find
 - De-clutters package space
 - No unique or long names for these types
 */
func (s *server) handleGreet() http.HandlerFunc {
	type request struct {
		Name string
	}
	type response struct {
		Greeting string `json:"greeting"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Actually handle the request ...
	}
}

/*
Lazy setup with sync.Once
 - perform expensive setup when the handler is first hit to improve startup time
 - if the handler isn't called, the work is never done
 */
func (s *server) handleTemplate(files string) http.HandlerFunc {
	var (
		init sync.Once
		tpl *template.Template
		tplerr error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func() {
			tpl, tplerr = template.ParseFiles(files)
		})
		if tplerr != nil {
			http.Error(w, tplerr.Error(), http.StatusInternalServerError)
			return
		}
		// use tpl
	}
}

/*
Middleware is just  Go functions
  - can run code before/after the wrapped handler
 - or choose not to call the wrapped handler at all
 */
func (s *server)  adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !currentUser(r).IsAdmin {
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}

/*
Respond helper
 - Abstract responding and do the bare bones initially
 - Later can make this more sophisticated (if needed)
 */
func (s *server) respond(w http.ResponseWriter, r *http.Request,  data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		err :=  json.NewEncoder(w).Encode(data)
		if err != nil {
			// Handle error
		}
	}
}

/*
Future proof helpers
 - Always take in http.ResponseWriter and *http.Request
 */

/*
Decoding helper - "Not a bad idea" - David Hernandez =  probably not a great idea (no no nope no)
 - Bare bones initially
 - can make more sophisticated (if needed)
 - could check content type later if needed
 */
func (s *server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func prepareThing() string {
	return ""
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	db, dbtidy, err := setupDatabase()
	if err != nil {
		return errors.Wrap(err, "setup database")
	}
	defer dbtidy()
	stv := &server{
		db: db,
	}
	return err
	// ... More stuff
}

func setupDatabase() (string, func(), error) {
	return "", func() {
		fmt.Println("Closing/tidying database")
	}, nil
}


