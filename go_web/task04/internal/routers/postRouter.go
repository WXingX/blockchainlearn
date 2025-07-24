package routers

import (
	"blog-management/internal/controllers"
	"github.com/gin-gonic/gin"
)

func PostRouterInit(r *gin.Engine) {
	postRouter := r.Group("/post")
	{
		postRouter.POST("/list", controllers.PostController{}.List)

		postRouter.POST("/create", controllers.PostController{}.Create)

		postRouter.GET("/list/detail/:post_id", controllers.PostController{}.Detail)

		postRouter.POST("/update", controllers.PostController{}.Edit)

		postRouter.POST("/delete", controllers.PostController{}.Delete)
	}
}
