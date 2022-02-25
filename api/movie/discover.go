package movie

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Discover(c *gin.Context) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.themoviedb.org/3/discover/movie", nil)
	if err != nil {
		log.Fatal(err)
	}

	// 设置参数
	params := req.URL.Query()

	params.Add("api_key", viper.GetString("themoviedb.apiKey"))
	params.Add("language", "zh")
	params.Add("sort_by", "popularity.desc")
	params.Add("with_watch_monetization_types", "flatrate")
	params.Add("page", "1")

	req.URL.RawQuery = params.Encode()

	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	var data interface{}
	json.Unmarshal(body, &data)

	c.JSON(http.StatusOK, data)
}
