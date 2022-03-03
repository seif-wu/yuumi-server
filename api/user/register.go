package user

import (
	"yuumi/internal/pkg/response"
	"yuumi/model"

	"github.com/gin-gonic/gin"
	v10 "github.com/go-playground/validator/v10"
)

func (h *handler) Register(c *gin.Context) {
	h.Lock()
	defer h.Unlock()
	// TODO 验证邮箱验证码是否正确

	var usernameRegisterParam UsernameRegisterParam
	if err := c.ShouldBindJSON(&usernameRegisterParam); err != nil {
		errs, ok := err.(v10.ValidationErrors)
		if !ok {
			response.ParamErr(c)
			return
		}
		response.Tran(c, "1010001", errs)
		return
	}

	var user model.User

	// 查询用户名是否存在
	h.DB.Select("Username").Where("username = ?", usernameRegisterParam.Username).First(&user)
	if user.Username != "" {
		response.Error(c, "1010002", "用户已存在")
		return
	}

	user = model.User{Username: usernameRegisterParam.Username, Password: usernameRegisterParam.Password}

	// 查询邀请码是否存在
	h.DB.Select("invitation_code").Where("invitation_code = ?", usernameRegisterParam.InviteCode).First(&user)
	if user.InvitationCode == "" {
		user.InvitationCode = ""
	}

	result := h.DB.Create(&user) // 通过数据的指针来创建
	if result.Error != nil {
		response.Error(c, "1010002", "用户创建失败")
		return
	}

	response.Success(c, "注册成功", gin.H{"data": user})
}
