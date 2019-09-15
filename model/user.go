package model

import (
	"fmt"
	"time"
)

type User struct {
	ID        string    `json:"id" validate:"required"`
	Firstname string    `json:"firstname,omitempty"`
	Lastname  string    `json:"lastname,omitempty"`
	Email     string    `json:"email,omitempty"`
	Username  string    `json:"username,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"` // TODO: Use pointer time
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Token     string    `json:"token,omitempty"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.Firstname, u.Lastname)
}

func GetUserTableName() string {
	return (&User{}).TableName()
}
