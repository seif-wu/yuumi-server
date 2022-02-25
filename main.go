package main

import (
	"fmt"
	"os"
	"yuumi/api/movie"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"https://api.themoviedb.org"})

	rootedPath, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(rootedPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	v1 := router.Group("/api/v1")
	{
		movies := v1.Group("movies")
		{
			movies.GET("/discover", movie.Discover)
		}
	}

	router.Run(":8088")
}
