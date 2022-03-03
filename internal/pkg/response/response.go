package response

import (
	"fmt"
	"net/http"
	"strconv"
	"yuumi/internal/pkg/validator"

	v10 "github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

/*
统一响应格式
{
	"success": bool,
	"data": map[string]interface{}
	"message": string,
	"code": string,
}

列表 success
	data []interface{}
	pageSize
	current
	total
*/

func ParamErr(c *gin.Context) {
	err := fmt.Errorf("参数绑定错误，请检查参数")
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    strconv.Itoa(http.StatusBadRequest),
		"message": err.Error(),
		"success": false,
	})
}

func Tran(c *gin.Context, code string, err v10.ValidationErrors) {
	Error(c, code, validator.TranslateStr(err, "\n"))
}

func Success(c *gin.Context, message string, result map[string]interface{}) {
	r := gin.H{
		"code":    http.StatusOK,
		"message": message,
		"success": true,
	}
	for key, value := range result {
		r[key] = value
	}

	c.JSON(http.StatusOK, r)
}

func Error(c *gin.Context, code string, err string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": err,
		"success": false,
	})
}
