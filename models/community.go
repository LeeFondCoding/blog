package models

import (
	"time"
)

type Community struct {
	CommunityID   uint64 `json:"community_id" db:"community_id"`
	CommunityName string `json:"community_name" db:"community_name"`
}

type CommunityDetail struct {
	CommunityID   uint64    `json:"community_id" db:"community_id"`
	CommunityName string    `json:"community_name" db:"community_name"`
	Introduction  string    `json:"introduction" db:"introduction"`
	CreateTime    time.Time `json:"create_time" db:"create_time"`
}

type CommunityDetailRes struct {
	CommunityID   uint64 `json:"community_id" db:"community_id"`
	CommunityName string `json:"community_name" db:"community_name"`
	Introduction  string `json:"introduction" db:"introduction"`
	CreateTime    string `json:"create_time" db:"create_time"`
}
