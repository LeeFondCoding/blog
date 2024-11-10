package redis

import (
	"errors"
)

var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
	ErrorVoteRepeat     = errors.New("不允许重复投票")
	ErrorVoted  = errors.New("已经投过票了"))