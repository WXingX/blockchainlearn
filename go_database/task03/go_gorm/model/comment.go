package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string
	UserID  uint
	PostID  uint
}

func (c *Comment) AfterDelete(db *gorm.DB) (err error) {
	//删除的时候 只给了一个id,所以c中只有id有值，其他都是零值，所以要先查出postID
	//可以通过 下面的方式查询，或者通过  db.First(&c, id) 来查询
	var postID uint
	fmt.Println("id:", c.ID)
	if c.PostID == 0 {
		err = db.Unscoped().Model(&Comment{}).
			Select("post_id").Where("id = ?", c.ID).Scan(&postID).Error
		if err != nil {
			fmt.Printf("AfterDelete 查询postid失败！%s \n", err.Error())
			return
		}
	}
	fmt.Printf("post id = %d \n", postID)
	var count int64
	err = db.Model(&Comment{}).Where("post_id = ?  AND deleted_at IS NULL", c.PostID).Count(&count).Error
	if err != nil {
		fmt.Println("AfterDelete failed...")
		return err
	}
	fmt.Printf("count = %d \n", count)
	if count == 0 {
		err = db.Model(&Post{}).
			Where("id = ?", postID).
			UpdateColumn("comment_status", "无评论").Error
	}

	return
}
