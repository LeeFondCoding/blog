package redis

import (
	"blog/models"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

// 按照分数从大到小的顺序查询指定数量的元素
func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	return client.ZRevRange(key, start, end).Result()
}

// 按order获取
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	key := KeyPostTimeZSet
	if p.Order == models.OrderScore {
		key = KeyPostScoreZSet
	}
	return getIDsFormKey(key, p.Page, p.Size)
}

func GetPostVoteData(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))
	for _, id := range ids {
		key := KeyPostVotedZSetPrefix + id
		// 查询key中分数是1的元素的数量
		v := client.ZCount(key, "1", "1").Val()
		data = append(data, v)
	}
	return data, nil
}

func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := KeyPostTimeZSet
	if p.Order == models.OrderScore {
		orderKey = KeyPostScoreZSet
	}
	communityKey := KeyCommunityPostSetPrefix + strconv.Itoa(int(p.CommunityID))
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		pipeline := client.Pipeline()
		
	}
}
