package logic

import (
	"blog/dao/api"
	"blog/models"
)

func GetGithubTrending(p *models.ParamGithubTrending) (
	data *models.GithubTrending, err error) {
	switch p.Language {
	case 0:
		data, err = api.GithubTrendingAll(p)
	}
	return
}
