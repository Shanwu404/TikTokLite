package service

import (
	"log"
	"strconv"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/middleware/redis"
	"github.com/Shanwu404/TikTokLite/utils"
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
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)

	key_userId := utils.Like_User_Key + strconv.FormatInt(userId, 10)
	key_videoId := utils.Like_Video_key + strconv.FormatInt(videoId, 10)

	//判断缓存中是否存在这个userid
	//如果存在，就补充一个键值对
	if n, err := redis.RDb.Exists(redis.Ctx, key_userId).Result(); n > 0 {
		if err != nil {
			log.Println("Redis query failed")
			return err
		}
		if _, err := redis.RDb.SAdd(redis.Ctx, key_userId, strVideoId).Result(); err != nil {
			log.Println("Redis add failed")
			return err
		} else { //添加进入数据库
			err := dao.InsertLike(&dao.Like{UserId: userId, VideoId: videoId})
			if err != nil {
				log.Println("the like operation error:", err.Error())
				return err
			}
			log.Println("Like operation successfully!")
			return nil
		}
	} else { //如果不存在这个user，就新建空白键值对
		redis.RDb.SAdd(redis.Ctx, key_userId, utils.MyDefault).Result()
		//设置过期时间
		redis.RDb.Expire(redis.Ctx, key_userId, redis.RandomTime()).Result()

		//把数据库中的当前用户点赞的videoId全部添加到缓存中
		videoIdList, err1 := dao.GetLikeVideoIdList(userId)
		if err1 != nil {
			log.Println("Failed to get the likes video id list")
			redis.RDb.Del(redis.Ctx, key_userId)
			return err1
		}
		for _, videoId := range videoIdList {
			strvideoId := strconv.FormatInt(videoId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RDb.SAdd(redis.Ctx, key_userId, strvideoId).Result(); err != nil {
				log.Println("Failed to add cache for videoId")
				redis.RDb.Del(redis.Ctx, key_userId)
				return err
			}
		}
		//把该次点赞的videoId添加到缓存中
		if _, err := redis.RDb.SAdd(redis.Ctx, key_userId, strVideoId).Result(); err != nil {
			log.Println("Failed to add cache for videoId")
			redis.RDb.Del(redis.Ctx, key_userId)
			return err
		} else { //只有缓存中添加正确，才可以往数据库中添加
			err := dao.InsertLike(&dao.Like{UserId: userId, VideoId: videoId})
			if err != nil {
				log.Println("the like operation error:", err.Error())
				return err
			}
			log.Println("Like operation successfully!")
		}
	}
	//判断缓存中是否存在这个videoid
	if n, err := redis.RDb.Exists(redis.Ctx, key_videoId).Result(); n > 0 {
		if err != nil {
			log.Println("Redis query failed")
			return err
		}
		//如果在缓存中，则把当前点赞的strUserId添加到key为strVideoId的set中
		if _, err := redis.RDb.SAdd(redis.Ctx, key_videoId, strUserId).Result(); err != nil {
			log.Println("Redis add failed")
			return err
		}
	} else {
		if _, err := redis.RDb.SAdd(redis.Ctx, key_videoId, utils.MyDefault).Result(); err != nil {
			log.Println("Cache creation key failed!")
			redis.RDb.Del(redis.Ctx, key_videoId)
			return err
		}
		if _, err := redis.RDb.Expire(redis.Ctx, key_videoId, redis.RandomTime()).Result(); err != nil {
			log.Println("Failed to set cache expiration time")
			redis.RDb.Del(redis.Ctx, key_videoId)
			return err
		}
		//把数据库中给当前视频的点赞的userId全部添加到缓存中
		userIdList, err1 := dao.GetLikeUserIdList(videoId)
		if err1 != nil {
			log.Println("Failed to get video id like user list")
			redis.RDb.Del(redis.Ctx, key_videoId)
			return err1
		}
		for _, userId := range userIdList {
			struserid := strconv.FormatInt(userId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RDb.SAdd(redis.Ctx, key_videoId, struserid).Result(); err != nil {
				log.Println("Failed to add cache for videoId")
				redis.RDb.Del(redis.Ctx, key_videoId)
				return err
			}
		}
		// if _, err := redis.RDb.SAdd(redis.Ctx, key_videoId, strUserId).Result(); err != nil {
		// 	log.Println("Failed to add cache for videoId")
		// 	redis.RDb.Del(redis.Ctx, key_videoId)
		// 	return err
		// }
	}
	return nil

}

/*取消点赞*/
func (like *LikeServiceImpl) Unlike(userId int64, videoId int64) error {
	err := dao.DeleteLike(userId, videoId)
	if err != nil {
		log.Println("the unlike operation error:", err.Error())
		return err
	}
	log.Println("Unlike operation successfully!")
	return nil
}

/*获取点赞列表, 返回的是视频的详细信息*/
func (like *LikeServiceImpl) GetLikeLists(userId int64) []VideoParams {
	videos, _ := dao.GetLikeVideoIdList(userId)
	results := make([]VideoParams, 0, len(videos))
	for _, video := range videos {
		log.Println(int64(video))
		result := like.videoService.QueryVideoById(int64(video))
		results = append(results, VideoParams(result))
	}
	log.Println("like list getting successfully!")
	//log.Println(like.videoService.QueryVideoById(1))
	return results
}

/*增加视频videoId的点赞数*/
func (like *LikeServiceImpl) addVideoLikeCount(videoId int64, sum *int64) {
	count, err := dao.CountLikes(videoId)
	if err != nil {
		log.Println("video likes adding failed")
		return
	}
	log.Println("the number of like getting successfully!")
	*sum += count
}

/*获取用户userId喜欢的视频数量*/
func (like *LikeServiceImpl) LikeVideoCount(userId int64) (int64, error) {
	likevideoIdList, err1 := dao.GetLikeVideoIdList(userId)
	if err1 != nil {
		log.Println("Failed to get the likes video id list")
		return 0, err1
	}
	log.Println("the number of like getting successfully!")
	return int64(len(likevideoIdList)), nil

}

/*判断用户userId是否点赞视频videoId*/
func (like *LikeServiceImpl) IsLike(videoId int64, userId int64) bool {
	videoIdList, err := dao.GetLikeVideoIdList(userId)
	if err != nil {
		log.Println("Failed to get the likes video id list")
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
		log.Println("count from db error:", err)
		return 0
	}
	log.Println("count likes successfully!")
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
			log.Println("Error counting likes for video ID:", video.ID, "Error:", err)
			continue // 如果发生错误, 记录错误并继续处理下一个视频
		}
		totalFavorites += likesForVideo
	}

	return totalFavorites
}
