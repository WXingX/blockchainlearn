package routers

import (
	"blog-management/internal/controllers"
	"github.com/gin-gonic/gin"
)

func UserRouterInit(r *gin.Engine) {
	userRouter := r.Group("/user")
	{
		//用户登陆
		userRouter.POST("/login", controllers.UserController{}.Login)
		//用户注册
		userRouter.POST("/register", controllers.UserController{}.Register)
	}
}
