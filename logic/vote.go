package logic

import (
	"blog/dao/redis"
	"blog/models"

	"go.uber.org/zap"

	"strconv"
)

// 为帖子投票
func VoteForPost(userID uint64, p *models.VoteDataForm) error {
	zap.L().Debug("VoteForPost",
		zap.Uint64("userId", userID),
		zap.String("postId", p.PostID),
		zap.Int8("Direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
