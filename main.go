package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"yuumi/internal/pkg/config"
	"yuumi/internal/pkg/logger"
	"yuumi/internal/pkg/mysql"
	"yuumi/internal/pkg/redis"
	apiRouter "yuumi/router"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// Gin 配置
	// 强制日志颜色化
	gin.ForceConsoleColor()
	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	// 加载配置
	if err := config.Init(); err != nil {
		panic("初始化配置错误")
	}
	// 初始化日志
	if err := logger.Init(); err != nil {
		panic("初始化日志错误")
	}
	// 初始化 Redis
	redisOption := redis.Option{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	RDB := redis.Init(&redisOption)
	// 初始化 MySQL
	mysqlOption := mysql.Options{
		Host:                  viper.GetString("mysql.host"),
		Username:              viper.GetString("mysql.user"),
		Password:              viper.GetString("mysql.password"),
		Database:              viper.GetString("mysql.dbname"),
		MaxIdleConnections:    200,
		MaxOpenConnections:    200,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
	}
	DB, err := mysql.Init(&mysqlOption)
	if err != nil {
		panic(fmt.Errorf("初始化MySQL错误： %s", err.Error()))
	}

	// 安装路由
	apiRouter.InstallRouter(router, DB, RDB)

	// 创建 server
	srv := &http.Server{
		Addr:    ":" + viper.GetString("app.port"),
		Handler: router,
	}
	// 启动 server
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
