package models

import (
	"time"
)

type User struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	Username     string `json:"username" gorm:"unique;not null"`
	Password     string `gorm:"not null" json:"password"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	IsAdmin      bool   `json:"isAdmin" gorm:"default:false"`
	Avatar       string `json:"avatar"`
	CoverImage   string `json:"coverImg"`
	Description  string `json:"desc"`
	Address      string `json:"address"`
	WorkAt       string `json:"workAt"`
	Relationship string `json:"relationship"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Follower struct {
	ID           int  `json:"id" gorm:"primaryKey"`
	UserId       int  `json:"userId"`
	FollowerId   int  `json:"followerId"`
	UserInfo     User `gorm:"foreignKey:UserId"`
	FollowerInfo User `gorm:"foreignKey:FollowerId"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
