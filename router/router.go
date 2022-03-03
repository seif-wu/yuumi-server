package router

import (
	"fmt"
	"yuumi/api/auth"
	"yuumi/api/movie"
	"yuumi/api/user"
	"yuumi/model"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func InstallRouter(router *gin.Engine, db *gorm.DB, rdb *redis.Client) {
	if err := AutoMigrate(db); err != nil {
		panic(fmt.Errorf("migrate error，%s", err.Error()))
	}

	authMiddleware, err := auth.AuthMiddleware(db)
	if err != nil {
		panic(fmt.Errorf("auth error，%s", err))
	}

	v1 := router.Group("/api/v1")
	{
		v1.POST("/login", authMiddleware.LoginHandler)
		v1.POST("/logout", authMiddleware.LogoutHandler)
		v1.POST("/refresh_token", authMiddleware.RefreshHandler)
	}

	userController := user.NewController(db)

	v1Manager := router.Group("/api/v1/manager")
	v1Manager.Use(authMiddleware.MiddlewareFunc())
	{
		// 用户
		userRouter := v1Manager.Group("user")
		{
			userRouter.GET("/current", userController.Current)
		}
	}

	v1Public := router.Group("/api/v1/public")
	{
		userRouter := v1Public.Group("user")
		{
			userRouter.POST("/register", userController.Register)
		}

		movieController := movie.NewController(db)
		movieRouter := v1Public.Group("movie")
		{
			movieRouter.GET("/discover", movieController.Discover)
		}
	}
}

func AutoMigrate(db *gorm.DB) error {
	err := db.Set("", "CONVERT TO CHARACTER SET utf8mb4 ").AutoMigrate(
		&model.User{},
	)

	if err != nil {
		return err
	}

	return nil
}
