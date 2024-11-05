package mysql

import (
	"blog/models"
	"database/sql"
	"errors"
	"github.com/jmoix/sqlx"
	"go.uber.org/zap"
	"strings"
)

// 获取数据库的帖子总数
func GetPostTotalCount() (count int64, err error) {
	sqlStr := "select count(post_id) from post"
	err = db.Get(&count, sqlStr)
	if err != nil {
		zap.L().Error("db.Get failed", zap.Error(err))
		return 0, err
	}
	return
}

// 根据社区id，获取该社区下的帖子总数
func GetCommunityPostTotalCount(communityID int64) (count int64, err error) {
	sqlStr := "select count(post_id) from post where community_id = ?"
	err = db.Get(&count, sqlStr, communityID)
	if err != nil {
		zap.L().Error("db.Get failed", zap.Error(err))
		return 0, err
	}
	return
}

// 创建帖子
func CreatePost(post *models.Post) (err error) {
	sqlStr := `insert into post(
	post_id, title, content, author_id, community_id)
	values(?,?,?,?,?)`

	_, err = db.Exec(sqlStr, post.PostID, post.Title, post.Content, post.AuthorId, post.CommunityID)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

// 根据id查询帖子数据
func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id = ?`

	err = db.Get(post, sqlStr, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(ErrorInvalidID)
		}
		zap.L().Error("query post failed", zap.String("sql", sqlStr), zap.Error(err))
		return nil, errors.New(ErrorQueryFailed)
	}
	return
}

// 给定id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.ApiPostDetail, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)`

	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}

func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	ORDER BY create_time
	DESC 
	limit ?,?
	`

	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}
