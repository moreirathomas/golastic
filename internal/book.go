package internal

import "time"

type Book struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`

	Title    string `json:"title"`
	Abstract string `json:"abstract"`
	Author   Author `json:"author"`
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type BookService interface{}
