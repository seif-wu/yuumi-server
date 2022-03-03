package scope

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 分页 Scope
func Paginate(c *gin.Context) (scope func(db *gorm.DB) *gorm.DB, current int, pageSize int) {
	current, _ = strconv.Atoi(c.DefaultQuery("current", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 20
	}

	offset := (current - 1) * pageSize
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(pageSize)
	}, current, pageSize
}
