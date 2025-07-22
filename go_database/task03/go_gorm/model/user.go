package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string
	PassWord string
	Email    string
	PostNum  int
	Posts    []Post
}
