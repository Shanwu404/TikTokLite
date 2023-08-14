package service

import (
	"log"

	"github.com/Shanwu404/TikTokLite/dao"
)

type LikeServiceImpl struct {
	UserService
	VideoService
}

/*点赞*/
func (like LikeServiceImpl) Like(userId int64, videoId int64) error {
	err := dao.InsertLike(&dao.Like{UserId: userId, VideoId: videoId})
	if err != nil {
		log.Println("the like operation error:", err.Error())
		return err
	}
	log.Println("Like operation successfully!")
	return nil
}

/*取消点赞*/
func (like LikeServiceImpl) Unlike(userId int64, videoId int64) error {
	err := dao.DeleteLike(userId, videoId)
	if err != nil {
		log.Println("the unlike operation error:", err.Error())
		return err
	}
	log.Println("Unlike operation successfully!")
	return nil
}

/*获取点赞列表, 目前只有id，以后要update成返回详细信息*/
//返回[]dao.VideoDetail
func (like LikeServiceImpl) GetLikeLists(userId int64) ([]int64, error) {
	var videoids []int64
	videoids, err := dao.GetLikeVideoIdList(userId)
	if err != nil {
		log.Println("the like list getting error:", err.Error())
		return videoids, err
	}
	log.Println("like list getting successfully!")
	return videoids, nil

}

/*获取视频videoId的点赞数*/
func (like LikeServiceImpl) LikeCount(videoId int64) (int64, error) {
	var likeUserIdList []int64
	var result int64 = 0
	likeUserIdList, err := dao.GetLikeUserIdList(videoId)
	result = int64(len(likeUserIdList))
	if err != nil {
		log.Println("the number of like getting error:", err.Error())
		return result, err
	}
	log.Println("the number of like getting successfully!")

	return result, nil
}

/*增加视频videoId的点赞数*/
func (like LikeServiceImpl) addVideoLikeCount(videoId int64, sum *int64) {
	count, err := like.LikeCount(videoId)
	if err != nil {
		log.Println("video likes adding failed")
		return
	}
	log.Println("the number of like getting successfully!")
	*sum += count
}

/*获取用户userId喜欢的视频数量*/
func (like LikeServiceImpl) LikeVideoCount(userId int64) (int64, error) {
	likevideoIdList, err1 := dao.GetLikeVideoIdList(userId)
	if err1 != nil {
		log.Println("Failed to get the likes video id list")
		return 0, err1
	}
	log.Println("the number of like getting successfully!")
	return int64(len(likevideoIdList)), nil

}

/*判断用户userId是否点赞视频videoId*/
func (like LikeServiceImpl) IsLike(videoId int64, userId int64) (bool, error) {
	videoIdList, err := dao.GetLikeVideoIdList(userId)
	if err != nil {
		log.Println("Failed to get the likes video id list")
		return false, err
	}
	for _, vId := range videoIdList {
		if vId == videoId {
			return true, nil
		}
	}
	return false, nil
}

/*获取用户userId的获取的点赞总数*/
//func (like LikeServiceImpl) TotalLiked(userId int64) int64 {
//videoIdList := VideoServiceImpl{}.GetVideoIdListByUserId(userId)
//listlLen := len(videoIdList)
//videoLikecountList := make([]int64,listlLen)