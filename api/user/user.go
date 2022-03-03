package user

import (
	"sync"

	"gorm.io/gorm"
)

type handler struct {
	sync.RWMutex
	DB *gorm.DB
}

func NewController(db *gorm.DB) *handler {
	return &handler{
		DB: db,
	}
}

type UsernameRegisterParam struct {
	Username string `json:"username" binding:"required" zh:"用户名" en:"username"`
	Password string `json:"password" binding:"required" zh:"密码" en:"password"`
	// VerificationCode string `json:"verification_code" binding:"required" zh:"验证码" en:"verification code"`
	InviteCode string `json:"inviteCode"`
}
