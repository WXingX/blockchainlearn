package controllers

import (
	"blog-management/database"
	"blog-management/database/model"
	"blog-management/internal/response"
	"blog-management/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type AddCommentReq struct {
	Content string `json:"content"`
	PostId  uint   `json:"post_id"`
}

type CommentController struct{}

func (com CommentController) Create(c *gin.Context) {
	// TODO 创建评论
	var req AddCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, -1, err.Error())
		return
	}

	userId := c.GetUint("userID")
	if userId == 0 {
		response.Fail(c, http.StatusBadRequest, -1, "token is invalid")
		return
	}

	comment := model.Comment{
		Content: req.Content,
		PostID:  req.PostId,
		UserID:  userId,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		logger.Logger.Error("create comment error",
			zap.Uint("user_id", userId),
			zap.Uint("post_id", req.PostId),
			zap.String("content", req.Content),
			zap.Error(err))

		response.Fail(c, http.StatusInternalServerError, -1, "create comment error")
		return
	}

	response.Success(c, nil)
}
