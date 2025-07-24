package controllers

import (
	"blog-management/database"
	"blog-management/database/model"
	"blog-management/internal/response"
	"blog-management/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserController struct{}

func (u UserController) Login(c *gin.Context) {
	//	TODO 用户登录
	var loginReq LoginReq
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		response.Fail(c, http.StatusBadRequest, -1, err.Error())
		return
	}

	var storedUser model.User
	if err := database.DB.Where("username = ?", loginReq.Username).First(&storedUser).Error; err != nil {
		response.Fail(c, http.StatusBadRequest, -1, "Invalid username or password")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(loginReq.Password)); err != nil {
		response.Fail(c, http.StatusUnauthorized, -1, "Invalid username or password")
		return
	}

	// 生成 JWT
	token, err := utils.GenToken(storedUser.ID, storedUser.Username)

	if err != nil {
		response.Fail(c, http.StatusInternalServerError, -1, "Failed to generate token")
		return
	}

	// 在header中添加token
	c.Header("Authorization", "Bearer "+token)
	response.Success(c, nil)
}

func (u UserController) Register(c *gin.Context) {
	//用户注册
	//注册信息JSON 包括  username  email  password
	var registerReq RegisterReq
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		response.Fail(c, http.StatusBadRequest, -1, err.Error())
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, -1, "Failed to hash password")
		return
	}
	registerReq.Password = string(hashedPassword)
	user := model.User{
		Username: registerReq.Username,
		Password: registerReq.Password,
		Email:    registerReq.Email,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, -1, "Failed to create user")
		return
	}

	response.Success(c, nil)
}
