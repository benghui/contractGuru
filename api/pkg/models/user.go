package models

import (
	"html"
	"strings"
)

// Users contain parameters for user
type User struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (u *User) Prepare() {
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
}
