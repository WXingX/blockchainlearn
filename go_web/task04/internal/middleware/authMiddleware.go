package middleware

import (
	"blog-management/internal/response"
	"blog-management/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"strings"
)

const (
	bearerWord       string = "Bearer"
	authorizationKey string = "Authorization"
)

var whiteList = []string{
	"/user/login",
	"/user/register",
}

func AuthJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 过滤一部分路由
		if slices.Contains(whiteList, c.FullPath()) {
			c.Next()
			return
		}
		authHeader := c.Request.Header.Get(authorizationKey)
		if authHeader == "" {
			response.Fail(c, http.StatusUnauthorized, -1, "token is missing")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == bearerWord) {
			response.Fail(c, http.StatusUnauthorized, -1, "token is invalid")
			c.Abort()
			return
		}

		mc, err := utils.ParseToken(parts[1])
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, -1, "token is invalid")
			c.Abort()
			return
		}
		c.Set("userID", mc.UserID)
		c.Set("userName", mc.Username)
		c.Next()
	}
}
