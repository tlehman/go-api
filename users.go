package main

import (
	"encoding/json"
	"fmt"

	"github.com/replit/database-go"
)

type User struct {
	name  string
	id    string
	roles string
}

func parseUsers(value string) []User {
	var users []User
	err := json.Unmarshal([]byte(value), &users)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return users
}

func registerUser(user User) {
	usersRaw, err := database.Get("users")
	var users []User = parseUsers(usersRaw)
	registeredUsers = users

	loggedInUserIsRegistered := false
	if err != nil {
		fmt.Println(err)
	} else {
		// check if logged in user is registered
		for _, user := range users {
			if user.id != "" && (user.id == loggedInUser.id) {
				loggedInUserIsRegistered = true
			}
		}
	}
	if !loggedInUserIsRegistered {
		users = append(users, loggedInUser)
		usersJson, err := json.Marshal(users)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		database.Set("users", string(usersJson))
	}
}
