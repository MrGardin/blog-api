package storage

import (
	"time"
)

type User struct {
	Id         int       `json:"id"`
	FirstName  string    `json:"first_name"`
	SecondName string    `json:"second_name"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
}
