package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title         string
	Content       string
	CommentStatus string
	UserID        uint
	Comments      []Comment
}

// AfterCreate 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
func (post *Post) AfterCreate(db *gorm.DB) (err error) {
	var count int64
	err = db.Model(&Post{}).Where("user_id = ?", post.UserID).
		Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.Model(&User{}).Where("id = ?", post.UserID).
		UpdateColumn("post_num", count).Error
	return
}
