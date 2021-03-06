package base

import (
	"time"

	"gorm.io/gorm"
)

type ModelBase struct {
	Id        int            `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}
