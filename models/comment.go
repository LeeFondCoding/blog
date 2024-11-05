package models

import (
	"time"
)

type Comment struct {
	PostID    uint64 `db:"post_id" json:"post_id"`
	ParentID  uint64 `db:"parent_id" json:"parent_id"`
	CommentID uint64 `db:"comment_id" json:"comment_id"`
	AuthorID  uint64 `db:"author_id" json:"author_id"`
	Content   string `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
