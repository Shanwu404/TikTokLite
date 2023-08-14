package service

import (
	"log"

	"github.com/Shanwu404/TikTokLite/dao"
)

type RelationServiceImpl struct{}

/*
关注用户
userId 关注 followId
*/
func (RelationServiceImpl) Follow(userId int64, followId int64) error {
	rsi := RelationServiceImpl{}
	isFollowed := rsi.JudgeIsFollowById(userId, followId)
	//如果没有关注， 那么进行关注操作
	if !isFollowed {
		err := dao.InsertFollow(userId, followId)
		if err != nil {
			log.Println("Error exists when doing follow operation")
			return err
		}
		log.Println("Success in follow operation")
	}
	log.Println("You have aleady followed him, unvalid follow operation")
	return nil
}

/*
取关用户
userId 关注 followId
*/
func (RelationServiceImpl) UnFollow(userId int64, followId int64) error {
	rsi := RelationServiceImpl{}
	isFollowed := rsi.JudgeIsFollowById(userId, followId)

	//如果关注了， 那么进行取关操作
	if isFollowed {
		err := dao.DeleteFollow(userId, followId)
		if err != nil {
			log.Println("Error exists when doing unfollow operation")
			return err
		}
		log.Println("Success in unfollow operation")
	}
	log.Println("You haven't followed him, unvalid unfollow operation")
	return nil
}

/*
查询是否已关注
userId 关注 followId
*/
func (RelationServiceImpl) JudgeIsFollowById(userId int64, followId int64) bool {
	flag := dao.JudgeIsFollowById(userId, followId)
	return flag
}

/*
获取用户关注列表
*/
func (RelationServiceImpl) GetFollowList(userId int64) ([]dao.UserResp, error) {
	usi := UserServiceImpl{}
	followList := make([]dao.UserResp, 0)
	followIds, err := dao.QueryFollowsIdByUserId(userId)
	if nil != err {
		return followList, err
	}
	for _, followId := range followIds {
		followInfo, err := usi.QueryUserRespByID(followId)
		if nil != err {
			return followList, err
		}
		followInfo.IsFollow = true
		followList = append(followList, followInfo)
	}
	return followList, nil
}

/*
获取用户粉丝列表
*/
func (RelationServiceImpl) GetFollowerList(userId int64) ([]dao.UserResp, error) {
	rsi := RelationServiceImpl{}
	usi := UserServiceImpl{}
	followerList := make([]dao.UserResp, 0)
	followerIds, err := dao.QueryFollowersIdByUserId(userId)
	if nil != err {
		return followerList, err
	}

	for _, followerId := range followerIds {
		followerInfo, err := usi.QueryUserRespByID(followerId)
		isFollow := rsi.JudgeIsFollowById(userId, followerId)
		if nil != err {
			return followerList, err
		}
		if isFollow {
			followerInfo.IsFollow = true
		} else {
			followerInfo.IsFollow = false
		}
		followerList = append(followerList, followerInfo)
	}
	return followerList, nil
}

/*
获取用户好友列表
*/
func (RelationServiceImpl) GetFriendList(userId int64) ([]dao.FriendResp, error) {
	friendList := make([]dao.FriendResp, 0)
	rsi := RelationServiceImpl{}

	// 查出关注列表
	usi := UserServiceImpl{}
	followIds, err := dao.QueryFollowsIdByUserId(userId)
	if nil != err {
		return friendList, err
	}
	for _, followId := range followIds {
		tmpFriendInfo, err := usi.QueryUserRespByID(followId)
		friendResp := dao.FriendResp{}
		// 判断是否回关，回关了即为好友
		isFollow := rsi.JudgeIsFollowById(followId, userId)
		if nil != err {
			return friendList, err
		}
		if isFollow {
			friendResp.Id = int64(tmpFriendInfo.Id)
			friendResp.Name = tmpFriendInfo.Username
			friendResp.FollowCount = tmpFriendInfo.FollowCount
			friendResp.FollowerCount = tmpFriendInfo.FollowerCount
			friendResp.IsFollow = true
			friendResp.Avatar = "http://chy"
			//friendResp.Avatar = "http://" + config.Url_addr + config.Url_Image_prefix + "male.png"
			friendResp.FavoriteCount = tmpFriendInfo.FavoriteCount
			friendResp.WorkCount = tmpFriendInfo.WorkCount
			friendResp.TotalFavorited = tmpFriendInfo.TotalFavorited

			//待开发
			//msi := MessageServiceImpl{}
			//latestMsg := msi.QueryMessagesByIds(userId, int64(tmpFriendInfo.Id))
			// 如果没有最新消息,可以做一个容错
			// if err != nil {
			// 	friendResp.Message = ""
			// 	friendResp.MsgType = 0
			// }
			//friendResp.Message = latestMsg.Content
			friendList = append(friendList, friendResp)
		}
	}

	return friendList, nil
}

// 统计id用户粉丝数
func (RelationServiceImpl) CountFollowers(id int64) int64 {
	cnt := dao.CountFollowers(id)
	return cnt
}

// 统计id用户关注数
func (RelationServiceImpl) CountFollowings(id int64) int64 {

	cnt := dao.CountFollowees(id)
	return cnt
}
