package movie

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yuumi/internal/pkg/themoviedb"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type DiscoverQuery struct {
	Page                       string `json:"page"`
	Language                   string `json:"language"`
	SortBy                     string `json:"sort_by"`
	WithWatchMonetizationTypes string `json:"with_watch_monetization_types"`
}

func NewDiscoverQuery(c *gin.Context) DiscoverQuery {
	var discoverQuery DiscoverQuery
	discoverQuery.Language = c.DefaultQuery("language", "zh")
	discoverQuery.Page = c.DefaultQuery("page", "1")
	discoverQuery.SortBy = c.DefaultQuery("sort_by", "popularity.desc")
	discoverQuery.WithWatchMonetizationTypes = c.DefaultQuery("with_watch_monetization_types", "flatrate")

	return discoverQuery
}

func (u *handler) Discover(c *gin.Context) {
	client := themoviedb.Client(themoviedb.ClientConfig{
		ApiKey: viper.GetString("themoviedb.apiKey"),
	})

	discoverQuery := NewDiscoverQuery(c)

	fmt.Println(discoverQuery)

	var params map[string]string
	discoverQueryMap, _ := json.Marshal(discoverQuery)
	json.Unmarshal(discoverQueryMap, &params)
	data, _ := client.Discover(params)

	c.JSON(http.StatusOK, data)
}
