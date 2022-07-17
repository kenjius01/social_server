package models

import "time"

type Post struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	UserId    int    `json:"userId"`
	Desc      string `json:"desc"`
	Image     string `json:"image"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Like struct {
	ID        int `json:"id" gorm:"primaryKey"`
	UserId    int `json:"userId"`
	PostId    int `json:"postId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
