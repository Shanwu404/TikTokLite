package dao

import (
	"github.com/Shanwu404/TikTokLite/log/logger"
)

// 定义用户信息结构体
type UserInfo struct {
	Id              int64  `json:"id"`               // 用户ID
	Username        string `json:"name"`             // 用户名
	FollowCount     int64  `json:"follow_count"`     // 关注数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝数
	IsFollow        bool   `json:"is_follow"`        // 是否关注
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 背景图片
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  int64  `json:"total_favorited"`  // 获赞数
	WorkCount       int64  `json:"work_count"`       // 作品数
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
}

// 定义关注关系结构体
type Follow struct {
	ToUserId   int64 `gorm:"column:to_user_id"`  // 被关注用户ID
	FollowerId int64 `gorm:"column:follower_id"` // 执行关注的用户ID
}

// InsertFollow 增加follow关系 userId 关注 followId
func InsertFollow(userId int64, followId int64) error {
	follow := Follow{
		ToUserId:   followId, // 被关注用户ID
		FollowerId: userId,   // 执行关注的用户ID
	}
	if err := db.Table("follows").Create(&follow).Error; err != nil {
		logger.Errorln(err)
		return err // 如果插入出错，则返回错误
	}
	logger.Infoln("Insert follow success")
	return nil
}

// DeleteFollow 删除follow关系 userId 关注 followId
func DeleteFollow(userId int64, followId int64) error {
	err := db.Table("follows").Where("to_user_id = ? and follower_id = ?", followId, userId).Delete(&Follow{}).Error
	if err != nil {
		logger.Errorln("Delete follow error: ", err.Error())
		return err
	}
	logger.Infoln("Delete follow success")
	return nil
}

// IsFollowed 查询是否已关注 userId 关注 followId
func IsFollowed(userId int64, followId int64) (bool, error) {
	var count int64
	err := db.Model(&Follow{}).Where("follower_id = ? and to_user_id = ?", userId, followId).Count(&count).Error
	if err != nil {
		logger.Errorln("Judge is follow error: ", err.Error())
		return false, err
	}
	logger.Infoln("Judge is follow by UserId success")
	return count > 0, nil
}

// QueryFollowersIdByUserId 查询用户粉丝id列表
func QueryFollowersIdByUserId(userId int64) ([]int64, error) {
	var followersId []int64
	err := db.Table("follows").Where("to_user_id = ?", userId).Pluck("follower_id", &followersId).Error
	if err != nil {
		logger.Errorln("Query followers error: ", err.Error())
		return nil, err
	}
	logger.Infoln("Query followers Id by UserId success")
	return followersId, nil
}

// QueryFollowsIdByUserId 查询用户关注id列表
func QueryFollowsIdByUserId(userId int64) ([]int64, error) {
	var followsId []int64
	err := db.Table("follows").Where("follower_id = ?", userId).Pluck("to_user_id", &followsId).Error
	if err != nil {
		logger.Errorln("Query follows error: ", err.Error())
		return nil, err
	}
	logger.Infoln("Query follows success")
	return followsId, nil
}

// CountFollowers 统计用户粉丝数
func CountFollowers(userId int64) (int64, error) {
	var count int64
	err := db.Table("follows").Where("to_user_id = ?", userId).Count(&count).Error
	if err != nil {
		logger.Errorln("Count followers error: ", err.Error())
		return 0, err
	}
	logger.Infoln("Count followers success")
	return count, nil
}

// CountFollows 统计用户关注数
func CountFollows(userId int64) (int64, error) {
	var count int64
	err := db.Table("follows").Where("follower_id = ?", userId).Count(&count).Error
	if err != nil {
		logger.Errorln("Count follows error: ", err.Error())
		return 0, err
	}
	logger.Infoln("Count follows success")
	return count, nil
}
