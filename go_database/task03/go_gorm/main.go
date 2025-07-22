package main

import (
	"fmt"
	"go_gorm/model"
	"net/url"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	db := InitDB()
	if db == nil {
		fmt.Println("main() InitDB failed.")
		os.Exit(0)
	}
	err := db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	if err != nil {
		fmt.Println("main() AutoMigrate failed.")
		os.Exit(0)
	}
	//addPost(db)
	//addComment(db)
	//SelectUser(db)
	//posts := GetAllPostAndCommentByUserId(db, 1)
	//for _, post := range posts {
	//	fmt.Println(post)
	//}
	//post := GetMaxCommentNumPost(db)
	//fmt.Println(post)
	deleteComment(db)
}

func InitDB() *gorm.DB {
	user := "xxxxx"
	password := "xxxxx"
	enPwd := url.QueryEscape(password)
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local", user, enPwd)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tbl_",
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Printf("InitDB failed, err: %s\n", err.Error())
		return nil
	}
	return db
}

func addUser(db *gorm.DB) {
	db.Create(&model.User{
		UserName: "abc",
		PassWord: "123",
		Email:    "abc@123.com",
	})
	db.Create(&model.User{
		UserName: "yyy",
		PassWord: "ssss",
		Email:    "yyy@ssss.com",
	})
	db.Create(&model.User{
		UserName: "xxx",
		PassWord: "www",
		Email:    "xxx@www.com",
	})
}

func addPost(db *gorm.DB) {
	//db.Create(&model.Post{
	//	Title:   "文章1",
	//	Content: "这是第一篇文章。。。",
	//	UserID:  1,
	//})
	//
	//db.Create(&model.Post{
	//	Title:   "文章2",
	//	Content: "这是第二篇文章。。。",
	//	UserID:  1,
	//})
	//
	//db.Create(&model.Post{
	//	Title:   "文章3",
	//	Content: "这是第三篇文章。。。",
	//	UserID:  2,
	//})
	//db.Create(&model.Post{
	//	Title:   "文章4",
	//	Content: "这是第五篇文章。。。",
	//	UserID:  1,
	//})
}

func addComment(db *gorm.DB) {
	//db.Create(&model.Comment{
	//	Content: "写的不错",
	//	UserID:  2,
	//	PostID:  1,
	//})
	//
	//db.Create(&model.Comment{
	//	Content: "写的勉强过的去",
	//	UserID:  3,
	//	PostID:  1,
	//})

	//db.Create(&model.Comment{
	//	Content: "不对，偏题了",
	//	UserID:  2,
	//	PostID:  2,
	//})
}

func deleteComment(db *gorm.DB) {
	var comment model.Comment = model.Comment{}
	comment.ID = 3
	db.Delete(&comment)
}

func SelectUser(db *gorm.DB) {
	result := make(map[string]interface{})
	db.Model(&model.User{}).First(&result)
	fmt.Println(result)
}

// GetAllPostAndCommentByUserId 使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
func GetAllPostAndCommentByUserId(db *gorm.DB, user_id uint) []model.Post {
	var posts []model.Post
	db.Preload("Comments").Where("user_id = ?", user_id).Find(&posts)
	return posts
}

// 编写Go代码，使用Gorm查询评论数量最多的文章信息。
func GetMaxCommentNumPost(db *gorm.DB) *model.Post {
	//SELECT post_id, COUNT(*) AS CNT FROM tbl_comment GROUP BY post_id ORDER BY cnt DESC LIMIT 1
	var post model.Post
	// 查询评论数量最多的文章 id
	result := make(map[string]interface{})
	db.Model(&model.Comment{}).
		Select("post_id, COUNT(*) AS cnt").
		Group("post_id").
		Order("cnt DESC").
		Limit(1).
		Scan(&result)
	fmt.Println(result)
	if len(result) > 0 {
		db.Model(&model.Post{}).First(&post, result["post_id"])
	} else {
		return nil
	}

	return &post
}
