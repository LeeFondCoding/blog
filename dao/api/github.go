package api

import (
	"blog/models"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const (
	// 获取github趋势榜的url
	oringUrl = "https://api.github.com/search/repositories?q=stars:%253E1&sort=stars&order=desc&page="
	perPage  = "&per_page="
)

// 获取github全语言热榜
func GithubTrendingAll(p *models.ParamGithubTrending) (data *models.GithubTrending, err error) {
	url := oringUrl + fmt.Sprintf("%d", p.Page) + perPage + fmt.Sprintf("%d", p.Size)

	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		zap.L().Error("http.NewRequest failed", zap.Error(err))
		return
	}

	res, err := client.Do(req)
	if err != nil {
		zap.L().Error("client.Do failed", zap.Error(err))
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		zap.L().Error("io.ReadAll failed", zap.Error(err))
		return
	}
	fmt.Println(string(body))

	var GithubTrendingAll models.GithubTrending
	err = json.Unmarshal(body, &GithubTrendingAll)
	if err != nil {
		zap.L().Error("json.Unmarshal failed", zap.Error(err))
		return
	}
	return &GithubTrendingAll, nil
}
