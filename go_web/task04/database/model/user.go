package model

type User struct {
	BaseModel
	Username string `gorm:"unique;not null" json:"user_name"`
	Password string `gorm:"not null" json:"password"`
	Email    string `gorm:"unique;not null" json:"email"`
}
