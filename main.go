package main

import (
	"fmt"
	"net/http"
	"html/template"
	"time"
	"path"
)

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
  
	data := struct{
		Time time.Time
    UserId string
    UserName string
    UserRoles string
	}{
		Time: time.Now(),
    UserId: userId,
    UserName: userName,
    UserRoles: userRoles,
	}

	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error", err)
	}
}

func main() {
	fmt.Println("http server up!")
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