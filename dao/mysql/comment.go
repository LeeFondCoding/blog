package mysql

import (
	"blog/models"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// 创建评论
func CreateComment(comment *models.Comment) (err error) {
	sqlStr := `insert into comment(
	comment_id, content, post_id, author_id, parent_id)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, comment.CommentID, comment.Content, comment.PostID, comment.AuthorID, comment.ParentID)
	if err != nil {
		zap.L().Error("CreateComment failed", zap.Error(err))
		return
	}
	return
}

// 查询给定id列表的评论列表
func GetCommentListByIDs(ids []string) (commentList []*models.Comment, err error) {
	sqlStr := `select comment_id, content, post_id, author_id, parent_id, create_time
	from comment
	where comment_id in (?)`
	query, args, err := sqlx.In(sqlStr, ids)
	if err != nil {
		zap.L().Error("sqlx.In failed", zap.Error(err))
		return
	}
	query = db.Rebind(query)
	err = db.Select(&commentList, query, args...)
	if err != nil {
		zap.L().Error("db.Select failed", zap.Error(err))
		return
	}
	return
}
