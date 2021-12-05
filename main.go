package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"time"
)

var loggedInUser User
var registeredUsers []User

func handle(w http.ResponseWriter, r *http.Request) {
	// You might want to move ParseGlob outside of handle so it doesn't
	// re-parse on every http request.
	tmpl, err := template.ParseGlob("templates/*")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := ""
	if r.URL.Path == "/" {
		name = "index.html"
	} else {
		name = path.Base(r.URL.Path)
	}

	// Get the authentication headers
	userId := r.Header.Get("X-Replit-User-Id")
	userName := r.Header.Get("X-Replit-User-Name")
	userRoles := r.Header.Get("X-Replit-User-Roles")

	isLoggedIn := false
	if userName != "" {
		loggedInUser = User{name: userName, id: userId, roles: userRoles}
		registerUser(loggedInUser)
		isLoggedIn = true
	}

	data := struct {
		Time            time.Time
		IsLoggedIn      bool
		LoggedInUser    User
		RegisteredUsers []User
	}{
		Time:            time.Now(),
		IsLoggedIn:      isLoggedIn,
		LoggedInUser:    loggedInUser,
		RegisteredUsers: registeredUsers,
	}

	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error", err)
	}
}

func main() {
	fmt.Println("go play go!")
	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static")),
		),
	)
	http.HandleFunc("/", handle)
	http.ListenAndServe(":0", nil)
}
