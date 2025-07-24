package routers

import (
	"blog-management/internal/controllers"
	"github.com/gin-gonic/gin"
)

func CommentRouterInit(r *gin.Engine) {
	CommentRouter := r.Group("/comment")
	{
		//实现评论的创建功能，已认证的用户可以对文章发表评论。
		CommentRouter.POST("/create", controllers.CommentController{}.Create)
	}
}
