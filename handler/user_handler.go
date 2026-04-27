package handler

import (
	"fmt"
	"go-base/config"
	"go-base/model"
	"go-base/response"
	"go-base/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	UserService *service.UserService
}

type RegisterRequest struct {
	// 用户姓名
	Username string `json:"username" binding:"required,min=2,max=50"`
	// 用户密码
	Password string `json:"password" binding:"required,min=6,max=50"`
	// 用户邮箱
	Email string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	// 用户姓名
	Username string `json:"username" binding:"required"`
	// 用户密码
	Password string `json:"password" binding:"required"`
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册参数"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err.Error())
		return
	}

	// 转换为 model.User
	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		response.Error(c, err.Error())
		return
	}

	// 返回成功响应
	message := fmt.Sprintf("User: %s registered successfully", user.Username)
	response.SuccessWithMessage(c, message, nil)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户名密码登录并返回 token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录参数"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err.Error())
		return
	}

	// 查询用户
	user, err := h.UserService.GetUserByUsername(req.Username)
	if err != nil {
		response.Error(c, "Invalid username or password")
		return
	}

	// 验证密码
	if err := h.UserService.VerifyPassword(user.Password, req.Password); err != nil {
		response.Error(c, "Invalid username or password")
		return
	}

	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(config.GlobalConfig.JWT.ExpireHour)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		response.Error(c, "Failed to generate token")
		return
	}

	// 返回成功响应
	response.Success(c, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// 获取当前登录用户信息
// @Summary 获取当前用户信息
// @Description 需要携带 Bearer Token
// @Tags 用户管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/userInfo [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 安全获取上下文中的 userId，避免类型断言 panic
	userId, exists := c.Get("userId")
	if !exists {
		response.Error(c, "Invalid user ID")
		return
	}

	user, err := h.UserService.GetUserById(userId)
	if err != nil {
		response.Error(c, "User not found")
		return
	}
	response.Success(c, user)
}
