package routers

import (
	"blog-management/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.TraceMiddleware(), middleware.LogMiddleware(), middleware.AuthJWTMiddleware())
	// 初始化路由
	UserRouterInit(r)
	PostRouterInit(r)
	CommentRouterInit(r)
	return r
}
