package model

type Comment struct {
	BaseModel
	Content string `gorm:"not null" json:"content"`
	UserID  uint   `json:"user_id"`
	User    User   `json:"-"`
	PostID  uint   `json:"post_id"`
	Post    Post   `json:"-"`
}
