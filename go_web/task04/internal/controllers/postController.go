package controllers

import (
	"blog-management/database"
	"blog-management/database/model"
	"blog-management/internal/response"
	"blog-management/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type CreatePostReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type EditPostReq struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type DeletePostReq struct {
	ID uint `json:"id"`
}

type PostController struct{}

func (p PostController) Create(c *gin.Context) {
	//	新增文章
	var req CreatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, -1, err.Error())
		return
	}
	userId := c.GetUint("userID")
	if userId == 0 {
		response.Fail(c, http.StatusBadRequest, -1, "token is invalid")
		return
	}
	post := model.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userId,
	}
	//	将post写入数据库中
	if err := database.DB.Create(&post).Error; err != nil {
		logger.Logger.Error("Failed to create post", zap.String("title", req.Title), zap.String("content", req.Content), zap.Uint("userId", userId), zap.Error(err))
		response.Fail(c, http.StatusInternalServerError, -1, "Failed to insert post")
		return
	}

	response.Success(c, nil)
}

func (p PostController) List(c *gin.Context) {
	// 获取当前用户的文章列表
	userId := c.GetUint("userID")
	if userId == 0 {
		response.Fail(c, http.StatusBadRequest, -1, "token is invalid")
		return
	}

	var posts []model.Post
	if err := database.DB.Where("user_id = ?", userId).Find(&posts).Error; err != nil {
		logger.Logger.Error("post list failed", zap.Uint("userId", userId), zap.Error(err))
		response.Fail(c, http.StatusInternalServerError, -1, "Failed to list posts")
		return
	}

	response.Success(c, posts)
}

func (p PostController) Detail(c *gin.Context) {
	// 获取文章详细信息 包括所有评论信息
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, -1, "postID is invalid")
		return
	}

	userId := c.GetUint("userID")
	if userId == 0 {
		response.Fail(c, http.StatusBadRequest, -1, "token is invalid")
		return
	}

	var post model.Post
	if err := database.DB.Where("id = ? AND user_id = ?", postID, userId).First(&post).Error; err != nil {
		logger.Logger.Error("get post detail failed", zap.Uint("UserId", userId), zap.Uint64("PostID", postID), zap.Error(err))
		response.Fail(c, http.StatusInternalServerError, -1, "Failed to get post detail")
		return
	}
	//获取评论并返回
	var comments []model.Comment
	if err := database.DB.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		logger.Logger.Error("get post comments failed", zap.Uint64("PostID", postID), zap.Error(err))
		response.Fail(c, http.StatusInternalServerError, -1, "Failed to get post comments")
		return
	}

	response.Success(c, struct {
		model.Post
		Comments []model.Comment `json:"comments"`
	}{
		post,
		comments,
	})
}

func (p PostController) Edit(c *gin.Context) {
	//	修改文章内容
	var req EditPostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Error("parse post request failed", zap.Error(err))
		response.Fail(c, http.StatusBadRequest, -1, err.Error())
		return
	}

	userId := c.GetUint("userID")
	if userId == 0 {
		response.Fail(c, http.StatusBadRequest, -1, "token is invalid")
		return
	}

	//	先查出 userid postid 对应的文章
	var post model.Post
	if err := database.DB.Where("id = ? AND user_id = ?", req.ID, userId).First(&post).Error; err != nil {
		logger.Logger.Error("edit post failed", zap.Uint("UserId", userId), zap.Uint("PostID", req.ID), zap.Error(err))
		response.Fail(c, http.StatusInternalServerError, -1, "Failed to edit post")
		return
	}
	if post.ID == 0 {
		response.Fail(c, http.StatusForbidden, -1, "auth forbidden")
		return
	}

	// 修改相关文章
	updateColumn := make(map[string]interface{})
	if req.Title != "" {
		updateColumn["title"] = req.Title
	}
	if req.Content != "" {
		updateColumn["content"] = req.Content
	}

	if err := database.DB.Model(&model.Post{}).Where("id = ?", req.ID).Updates(updateColumn).Error; err != nil {
		logger.Logger.Error("edit post failed",
			zap.Uint("PostID", req.ID),
			zap.String("title", req.Title),
			zap.String("content", req.Content),
			zap.Error(err))
		response.Fail(c, http.StatusInternalServerError, -1, "Failed to edit post")
		return
	}

	response.Success(c, nil)
}

func (p PostController) Delete(c *gin.Context) {
	//	TODO 删除文章内容
	var req DeletePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Error("parse post request failed", zap.Error(err))
		response.Fail(c, http.StatusBadRequest, -1, err.Error())
		return
	}

	userId := c.GetUint("userID")
	if userId == 0 {
		response.Fail(c, http.StatusBadRequest, -1, "token is invalid")
		return
	}

	var post model.Post
	if err := database.DB.Where("id = ? AND user_id = ?", req.ID, userId).First(&post).Error; err != nil {
		logger.Logger.Error("delete post failed",
			zap.Uint("UserId", userId),
			zap.Uint("postID", req.ID),
			zap.Error(err))
		response.Fail(c, http.StatusInternalServerError, -1, "Failed to delete post")
		return
	}

	if post.ID == 0 {
		response.Fail(c, http.StatusForbidden, -1, "auth forbidden")
		return
	}
	if err := database.DB.Model(&model.Post{}).Delete(&model.Post{}, req.ID).Error; err != nil {
		logger.Logger.Error("delete post failed", zap.Uint("PostID", req.ID), zap.Error(err))
		response.Fail(c, http.StatusInternalServerError, -1, "Failed to delete post")
		return
	}

	response.Success(c, nil)
}
