package service

import (
	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/log/logger"
)

type LikeServiceImpl struct {
	videoService    VideoService
	relationService RelationService
}

func NewLikeService() LikeService {
	return &LikeServiceImpl{
		videoService:    NewVideoService(),
		relationService: NewRelationService(),
	}
}

/*点赞*/
func (like *LikeServiceImpl) Like(userId int64, videoId int64) error {
	err := dao.InsertLike(&dao.Like{UserId: userId, VideoId: videoId})
	if err != nil {
		logger.Errorln("the like operation error:", err)
		return err
	}
	logger.Infoln("Like operation successfully!")
	return nil
}

/*取消点赞*/
func (like *LikeServiceImpl) Unlike(userId int64, videoId int64) error {
	err := dao.DeleteLike(userId, videoId)
	if err != nil {
		logger.Errorln("the unlike operation error:", err)
		return err
	}
	logger.Infoln("Unlike operation successfully!")

	return nil
}

/*获取点赞列表, 返回的是视频的详细信息*/
func (like *LikeServiceImpl) GetLikeLists(userId int64) []VideoParams {
	videos, _ := dao.GetLikeVideoIdList(userId)
	results := make([]VideoParams, 0, len(videos))
	for _, video := range videos {
		result := like.videoService.QueryVideoById(int64(video))
		results = append(results, VideoParams(result))
	}
	logger.Infoln("like list getting successfully!")
	return results
}

/*增加视频videoId的点赞数*/
func (like *LikeServiceImpl) addVideoLikeCount(videoId int64, sum *int64) {
	count, err := dao.CountLikes(videoId)
	if err != nil {
		logger.Errorln("video likes adding failed:", err)
		return
	}
	logger.Infoln("the number of like getting successfully!")
	*sum += count
}

/*获取用户userId喜欢的视频数量*/
func (like *LikeServiceImpl) LikeVideoCount(userId int64) (int64, error) {
	likevideoIdList, err := dao.GetLikeVideoIdList(userId)
	if err != nil {
		logger.Errorln("Failed to get the likes video id list:", err)
		return 0, err
	}
	logger.Infoln("the number of like getting successfully!")
	return int64(len(likevideoIdList)), nil

}

/*判断用户userId是否点赞视频videoId*/
func (like *LikeServiceImpl) IsLike(videoId int64, userId int64) bool {
	videoIdList, err := dao.GetLikeVideoIdList(userId)
	if err != nil {
		logger.Errorln("Failed to get the likes video id list:", err)
		return false
	}
	for _, vId := range videoIdList {
		if vId == videoId {
			return true
		}
	}
	return false
}

/*获取视频videoId的点赞数*/
func (like *LikeServiceImpl) CountLikes(videoId int64) int64 {
	cnt, err := dao.CountLikes(videoId)
	if err != nil {
		logger.Errorln("count from db error:", err)
		return 0
	}
	logger.Infoln("count likes successfully!")
	return cnt
}

/*获取用户userId发布视频的总被赞数*/
func (like *LikeServiceImpl) TotalFavorited(userId int64) int64 {
	// 获取该用户发布的所有视频
	videos := dao.QueryVideosByAuthorId(userId)

	totalFavorites := int64(0)

	// 遍历所有视频，获取每个视频的点赞数
	for _, video := range videos {
		likesForVideo, err := dao.CountLikes(video.ID) // 假设video有一个ID字段
		if err != nil {
			logger.Errorln("Error counting likes for video ID:", video.ID, "Error:", err)
			continue // 如果发生错误, 记录错误并继续处理下一个视频
		}
		totalFavorites += likesForVideo
	}

	return totalFavorites
}
