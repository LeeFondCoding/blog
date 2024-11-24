package logic

import (
	"blog/dao/mysql"
	"blog/models"
)

func GetCommunityList() ([]*models.Community, error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetailByID(id uint64) (*models.CommunityDetailRes, error) {
	return mysql.GetCommunityByID(id)
}