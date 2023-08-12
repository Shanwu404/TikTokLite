package dao

import (
	"log"
)

// 定义用户信息
type UserResp struct {
	Id              uint64 `json:"id"`               // 用户ID
	Username        string `json:"username"`         // 用户名
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

type Follow struct {
	UserId     int64 `gorm:"column:to_user_id"`  //被关注的人
	FollowerId int64 `gorm:"column:follower_id"` //点关注的人
}

type FriendResp struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	FollowCount    int64  `json:"follow_count"`
	FollowerCount  int64  `json:"follower_count"`
	IsFollow       bool   `json:"is_follow"`
	Avatar         string `json:"avatar"`
	TotalFavorited int64  `json:"total_favorited"`
	WorkCount      int64  `json:"work_count"`
	FavoriteCount  int64  `json:"favorite_count"`
	Message        string `json:"message"`
	//MsgType        int64  `json:"msgType"`
}

/*关注关系*/
// 增加关注关系，登录进来的用户（userid) 关注别人（followId)
func InsertFollow(userId int64, followId int64) error {
	follow := Follow{
		UserId:     followId,
		FollowerId: userId,
	}
	if err := db.Table("follows").Create(&follow).Error; err != nil {
		log.Println(err.Error())
		log.Println("the following operation fails")

		return err
	}
	return nil
}

// 删除关注关系 userId 不再关注 followId
func DeleteFollow(userId int64, followId int64) error {
	follow := Follow{}
	if err := db.Table("follows").Where("to_user_id = ? and follower_id = ?", followId, userId).Delete(&follow).Error; err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

/*关注列表*/
// JudgeIsFollowById 查询是否已关注 用户id1是否关注id2用户
func JudgeIsFollowById(id1 int64, id2 int64) bool { // 判断用户id1是否关注id2用户
	var count int64
	db.Model(&Follow{}).Where("to_user_id = ? and follower_id = ?", id2, id1).Count(&count)
	return count > 0
}

// QueryFollowsIdByUserId 查询用户关注id列表
func QueryFollowsIdByUserId(userId int64) ([]int64, error) {
	followIds := make([]int64, 0)
	if err := db.Table("follows").Select("to_user_id").Where("follower_id = ?", userId).Find(&followIds).Error; nil != err {
		return nil, err
	}
	return followIds, nil
}

/*粉丝列表*/
// QueryFollowersIdByUserId 查询用户粉丝id列表
func QueryFollowersIdByUserId(userId int64) ([]int64, error) {
	followerIds := make([]int64, 0)
	if err := db.Table("follows").Select("follower_id").Where("to_user_id = ?", userId).Find(&followerIds).Error; nil != err {
		return nil, err
	}
	return followerIds, nil
}

// CountFollowers 统计用户的粉丝数
func CountFollowers(id int64) int64 {
	var count int64
	db.Model(&Follow{}).Where("to_user_id = ?", id).Count(&count)
	return count
}

// CountFollowees 统计用户的关注数
func CountFollowees(id int64) int64 {
	var count int64
	db.Model(&Follow{}).Where("follower_id = ?", id).Count(&count)
	return count
}
