package themoviedb

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (client *ClientIns) Discover(movieDiscoverPrams map[string]string) (interface{}, error) {
	req, err := http.NewRequest("GET", "https://api.themoviedb.org/3/discover/movie", nil)
	if err != nil {
		return nil, err
	}

	// 设置参数
	params := req.URL.Query()
	params.Add("api_key", client.apiKey)

	for k, v := range movieDiscoverPrams {
		params.Add(k, v)
	}
	req.URL.RawQuery = params.Encode()

	resp, _ := client.httpClient.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	var data interface{}
	json.Unmarshal(body, &data)

	return data, nil
}
