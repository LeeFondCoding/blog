package logic

import (
	"blog/dao/mysql"
	"blog/dao/redis"
	"blog/models"
	"blog/pkg/snowflake"

	"fmt"
	"strconv"

	"go.uber.org/zap"
)

// 创建帖子
func CreatePost(post *models.Post) (err error) {
	postID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID() failed", zap.Error(err))
		return
	}
	post.PostID = postID

	if err := mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost(&post) failed", zap.Error(err))
		return err
	}

	community, err := mysql.GetCommunityNameByID(fmt.Sprint(post.CommunityID))
	if err != nil {
		zap.L().Error("mysql.GetCommunityNameByID failed", zap.Error(err))
		return err
	}

	if err := redis.CreatePost(
		postID,
		post.AuthorId,
		post.Title,
		TruncateByWords(post.Content, 120),
		community.CommunityID,
	); err != nil {
		return err
	}
	return
}

func GetPostById(postID int64) (data *models.ApiPostDetail, err error) {
	post, err := mysql.GetPostByID(postID)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(postID) failed",
			zap.Int64("postID", postID),
			zap.Error(err))
		return
	}

	// for ApiPostDetail.AuthorName
	user, err := mysql.GetUserByID(post.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserByID() failed",
			zap.Uint64("postID", post.AuthorId),
			zap.Error(err))
		return
	}

	community, err := mysql.GetCommunityByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed",
			zap.Uint64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}

	voteNum, err := redis.GetPostVoteNum(postID)
	data = &models.ApiPostDetail{
		Post:               post,
		CommunityDetailRes: community,
		AuthorName:         user.UserName,
		VoteNum:            voteNum,
	}
	return
}

func GetPostList(page, size int64) ([]*models.ApiPostDetail, error) {
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed",
			zap.Error(err))
		return nil, err
	}

	data := make([]*models.ApiPostDetail, 0, len(postList))
	for _, post := range postList {
		// for ApiPostDetail.AuthorName
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID", post.AuthorId),
				zap.Error(err))
			continue
		}

		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				zap.Uint64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		postDetail := &models.ApiPostDetail{
			Post:               post,
			CommunityDetailRes: community,
			AuthorName:         user.UserName,
		}
		data = append(data, postDetail)
	}
	return data, nil
}

func GetPostList2(p *models.ParamPostList) (*models.ApiPostDetailRes, error) {
	var res models.ApiPostDetailRes
	total, err := mysql.GetPostTotalCount()
	if err != nil {
		return nil, err
	}

	res.Page.Total = total
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder() return 0 data")
		return &res, nil
	}

	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return nil, err
	}

	res.Page.Page = p.Page
	res.Page.Size = p.Size
	res.List = make([]*models.ApiPostDetail, 0, len(posts))

	for idx, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID", post.AuthorId),
				zap.Error(err))
			user = nil
		}
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				zap.Uint64("community_id", post.CommunityID),
				zap.Error(err))
			community = nil
		}

		postDetail := &models.ApiPostDetail{
			VoteNum:            voteData[idx],
			Post:               post,
			CommunityDetailRes: community,
			AuthorName:         user.UserName,
		}
		res.List = append(res.List, postDetail)
	}
	return &res, nil
}

func GetCommunityPostList(p *models.ParamPostList) (*models.ApiPostDetailRes, error) {
	var res models.ApiPostDetailRes

	total, err := mysql.GetCommunityPostTotalCount(int64(p.CommunityID))
	if err != nil {
		return nil, err
	}
	res.Page.Total = total

	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIDsInOrder() return 0 data")
		return &res, nil
	}
	zap.L().Debug("GetCommunityPostList", zap.Any("ids", ids))

	vateData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return nil, err
	}

	res.Page.Page = p.Page
	res.Page.Size = p.Size
	res.List = make([]*models.ApiPostDetail, 0, len(posts))

	community, err := mysql.GetCommunityByID(p.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed",
			zap.Uint64("community_id", p.CommunityID),
			zap.Error(err))
		return nil, err
	}
	for idx, post := range posts {
		if post.CommunityID != p.CommunityID {
			continue
		}

		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID", post.AuthorId),
				zap.Error(err))
			user = nil
		}

		postDetail := &models.ApiPostDetail{
			VoteNum:            vateData[idx],
			Post:               post,
			CommunityDetailRes: community,
			AuthorName:         user.UserName,
		}
		res.List = append(res.List, postDetail)
	}
	return &res, nil
}

func GetPostListNew(p *models.ParamPostList) (data *models.ApiPostDetailRes, err error) {
	if p.CommunityID == 0 {
		//查所有
		data, err = GetPostList2(p)
	} else {
		//根据社区id查询
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew() failed", zap.Error(err))
		return nil, err
	}
	return
}

func PostSearch(p *models.ParamPostList) (*models.ApiPostDetailRes, error) {
	var res models.ApiPostDetailRes

	total, err := mysql.GetPostListTotoalCount(p)
	if err != nil {
		return nil, err
	}
	res.Page.Total = total

	posts, err := mysql.GetPostListByKeyWords(p)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return &models.ApiPostDetailRes{}, nil
	}

	ids := make([]string, 0, len(posts))
	for _, post := range posts {
		ids = append(ids, strconv.Itoa(int(post.PostID)))
	}

	vateData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	res.Page.Size = p.Size
	res.Page.Page = p.Page

	res.List = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID", post.AuthorId),
				zap.Error(err))
			user = nil
		}

		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				zap.Uint64("community_id", post.CommunityID),
				zap.Error(err))
			community = nil
		}

		postDetail := &models.ApiPostDetail{
			VoteNum:            vateData[idx],
			Post:               post,
			CommunityDetailRes: community,
			AuthorName:         user.UserName,
		}
		res.List = append(res.List, postDetail)
	}
	return &res, nil
}
