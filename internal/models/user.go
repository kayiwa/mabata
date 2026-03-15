package models

type User struct {
	Email  string
	Name   string
	Sub    string
	Groups []string
}
