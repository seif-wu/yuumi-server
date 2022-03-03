package model

import (
	"yuumi/model/base"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	base.ModelBase
	Username       string `json:"username" gorm:"size:64;comment:用户名"`
	Password       string `json:"-" gorm:"size:128;comment:密码"`
	NickName       string `json:"nick_name" gorm:"size:128;comment:昵称"`
	Mobile         string `json:"mobile" gorm:"size:11;comment:手机号"`
	Avatar         string `json:"avatar" gorm:"size:255;comment:头像"`
	Email          string `json:"email" gorm:"comment:邮箱"`
	IsAdmin        bool   `json:"-" gorm:"default:0;comment:是否为管理员"`
	InviteesID     int    `json:"invitees_id" gorm:"comment:被邀请人主键"`
	Invitees       *User  `json:"invitees"`
	InviteesCode   string `json:"invitees_code" gorm:"comment:被邀请人码"`
	InvitationCode string `json:"invitation_code" gorm:"comment:邀请码"`
}

func (User) TableName() string {
	return "users"
}

//加密
func (u *User) Encrypt() (err error) {
	if u.Password == "" {
		return
	}

	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		u.Password = string(hash)
		return
	}
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	return u.Encrypt()
}

func (u *User) BeforeUpdate(_ *gorm.DB) error {
	var err error
	if u.Password != "" {
		err = u.Encrypt()
	}
	return err
}
