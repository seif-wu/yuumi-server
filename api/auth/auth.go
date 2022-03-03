package auth

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"yuumi/internal/pkg/errcode"
	"yuumi/model"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var IdentityKey = "id"

func AuthMiddleware(db *gorm.DB) (authMiddleware *jwt.GinJWTMiddleware, err error) {
	// the jwt middleware
	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:            viper.GetString("jwt.Realm"),
		SigningAlgorithm: "HS256",
		Key:              []byte(viper.GetString("jwt.key")),
		Timeout:          2 * time.Hour,
		MaxRefresh:       720 * time.Hour,
		LoginResponse:    loginResponse(),
		LogoutResponse:   logoutResponse(),
		RefreshResponse:  refreshResponse(),
		IdentityKey:      IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			v, ok := data.(model.User)
			if ok {
				return jwt.MapClaims{
					IdentityKey: v.Id,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			var user model.User
			db.First(&user, claims[IdentityKey])

			return user
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginValues login
			if err := c.ShouldBind(&loginValues); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginValues.Username
			password := loginValues.Password

			var user model.User
			db.Where(&model.User{Username: username}).First(&user)

			if user.Id == 0 {
				return nil, jwt.ErrFailedAuthentication
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(model.User); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			zhMessage := message
			statusCode := code
			errorCode := strconv.Itoa(code)
			if message == "Token is expired" {
				zhMessage = "登录过期"
				errorCode = errcode.TokenExpiredCode
			}
			if message == "cookie token is empty" {
				zhMessage = "请先登录"
				errorCode = errcode.TokenEmpty
			}
			if message == "incorrect Username or Password" {
				zhMessage = "用户名或密码错误"
				statusCode = http.StatusUnauthorized
				errorCode = errcode.IncorrectUsernameOrPassword
			}

			c.JSON(statusCode, gin.H{
				"success": false,
				"code":    errorCode,
				"message": zhMessage,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	return authMiddleware, err
}

func loginResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"token":   token,
			"expire":  expire.Format(time.RFC3339),
		})
	}
}

func refreshResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"token":   token,
			"expire":  expire.Format(time.RFC3339),
		})
	}
}

func logoutResponse() func(c *gin.Context, code int) {
	return func(c *gin.Context, code int) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "退出登录成功",
		})
	}
}
