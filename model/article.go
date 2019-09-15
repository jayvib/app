package model

import "time"

type Article struct {
	ID        string    `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	Author    *Author   `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (a *Article) TableName() string {
	return "article"
}

func GetArticleTableName() string {
	return new(Article).TableName()
}
