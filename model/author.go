package model

import "time"

type Author struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (a *Author) TableName() string {
	return "author"
}

func GetAuthorTableName() string {
	return new(Author).TableName()
}
