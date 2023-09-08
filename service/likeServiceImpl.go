package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/log/logger"
	"github.com/Shanwu404/TikTokLite/middleware/rabbitmq"
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

	//判断是否已经进行点赞
	if isLike := like.IsLike(videoId, userId); isLike {
		return fmt.Errorf("user %d has already liked video %d", userId, videoId)
	}

	//若没有点赞则插入新的关系
	// if err := dao.InsertLike(&dao.Like{UserId: userId, VideoId: videoId}); err != nil {
	// 	return err
	// }
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)

	//若没有点赞则插入新的关系
	message := strings.Builder{}
	message.WriteString(strUserId)
	message.WriteString(":")
	message.WriteString(strVideoId)
	rabbitmq.RabbitMQLikeAdd.Producer(message.String())

	log.Println("insert success")

	key_userId := utils.LikeUserKey + strconv.FormatInt(userId, 10)
	key_videoId := utils.LikeVideokey + strconv.FormatInt(videoId, 10)

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
		}
	} else { //如果不存在这个user，就新建空白键值对
		redis.RDb.SAdd(redis.Ctx, key_userId, utils.MyDefault).Result()
		//设置过期时间
		redis.RDb.Expire(redis.Ctx, key_userId, utils.LikeUserKeyTTL).Result()

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
		if _, err := redis.RDb.Expire(redis.Ctx, key_videoId, utils.LikeVideoKeyTTL).Result(); err != nil {
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
	}
	return nil

}

/*取消点赞*/
func (like *LikeServiceImpl) Unlike(userId int64, videoId int64) error {
	//判断是否已经进行点赞
	if isLike := like.IsLike(videoId, userId); !isLike {
		return fmt.Errorf("Unlike operation error: user %d hasn't liked video %d", userId, videoId)
	}

	//如果已经点赞了，那么取消这次点赞
	// if err := dao.DeleteLike(userId, videoId); err != nil {
	// 	return err
	// }

	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)

	//如果已经点赞了，那么取消这次点赞
	message := strings.Builder{}
	message.WriteString(strUserId)
	message.WriteString(":")
	message.WriteString(strVideoId)
	rabbitmq.RabbitMQLikeDel.Producer(message.String())

	key_userId := utils.LikeUserKey + strconv.FormatInt(userId, 10)
	key_videoId := utils.LikeVideokey + strconv.FormatInt(videoId, 10)

	//查看用户是否在缓存中
	if n, err := redis.RDb.Exists(redis.Ctx, key_userId).Result(); n > 0 {
		if err != nil {
			log.Println("Redis query failed")
			return err
		} //如果存在缓存中，那么就将这个键值对删除
		if _, err := redis.RDb.SRem(redis.Ctx, key_userId, strVideoId).Result(); err != nil {
			return err
		}
	} else { //如果取消点赞的用户id不在redis缓存中,那么新建并设置过期时间
		redis.RDb.SAdd(redis.Ctx, key_userId, utils.MyDefault).Result()
		redis.RDb.Expire(redis.Ctx, key_userId, utils.LikeUserKeyTTL).Result()

		//把数据库中的当前用户点赞的videoId全部添加到缓存中
		videoIdList, err1 := dao.GetLikeVideoIdList(userId)
		if err1 != nil {
			return err1
		}
		for _, videoId := range videoIdList {
			strvideoId := strconv.FormatInt(videoId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RDb.SAdd(redis.Ctx, key_userId, strvideoId).Result(); err != nil {
				redis.RDb.Del(redis.Ctx, key_userId)
				return err
			}
		}
	}

	//查看该次取消点赞的strVideoId是否在缓存中
	if n, err := redis.RDb.Exists(redis.Ctx, key_videoId).Result(); n > 0 {
		if err != nil {
			log.Println("Redis query failed")
			return err
		}
		//如果在缓存中，则把当前点赞的strUserId移除key为strVideoId的set中
		if _, err := redis.RDb.SRem(redis.Ctx, key_videoId, strUserId).Result(); err != nil {
			log.Println("Redis removal failed")
			return err
		}
	} else { //如果不在缓存中,新建然后设置过期时间
		redis.RDb.SAdd(redis.Ctx, key_videoId, utils.MyDefault).Result()
		redis.RDb.Expire(redis.Ctx, key_videoId, utils.LikeVideoKeyTTL).Result()

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
	}

	return nil
}

/*获取点赞列表, 返回的是视频的详细信息*/
func (like *LikeServiceImpl) GetLikeLists(userId int64) []VideoParams {

	key_userId := utils.LikeUserKey + strconv.FormatInt(userId, 10)

	if n, err := redis.RDb.Exists(redis.Ctx, key_userId).Result(); n > 0 {
		if err != nil {
			log.Println("redis Query failed")
			return nil
		}
		videoIdList, err1 := redis.RDb.SMembers(redis.Ctx, key_userId).Result()
		if len(videoIdList) > 0 {
			videoIdList = videoIdList[1:]
		}
		log.Println(key_userId, videoIdList)
		if err1 != nil {
			log.Println("redis Query failed")
			return nil
		}

		videoIdListLen := len(videoIdList) //这个videoidlist在缓存中返回的应该是string类型的
		results := make([]VideoParams, 0, videoIdListLen)
		for _, id := range videoIdList {
			video, _ := strconv.ParseInt(id, 10, 64)
			result := like.videoService.QueryVideoById(int64(video))
			results = append(results, VideoParams(result))
		}
		log.Println("like list getting successfully!")
		return results
	} else { //如果key_userId不存在缓存中，需要把数据库中的信息添加到缓存中
		redis.RDb.SAdd(redis.Ctx, key_userId, utils.MyDefault).Result()
		redis.RDb.Expire(redis.Ctx, key_userId, utils.LikeUserKeyTTL).Result()

		//把数据库中的视频id添加到缓存中
		//把数据库中的当前用户点赞的videoId全部添加到缓存中
		videoIdList, err1 := dao.GetLikeVideoIdList(userId)
		if err1 != nil {
			log.Println("Failed to get the likes video id list")
			redis.RDb.Del(redis.Ctx, key_userId)
			return nil
		}
		for _, videoId := range videoIdList {
			strvideoId := strconv.FormatInt(videoId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RDb.SAdd(redis.Ctx, key_userId, strvideoId).Result(); err != nil {
				log.Println("Failed to add cache for videoId")
				redis.RDb.Del(redis.Ctx, key_userId)
				return nil
			}
		}

		videoIdListLen := len(videoIdList)
		results := make([]VideoParams, 0, videoIdListLen)
		for _, video := range videoIdList {
			result := like.videoService.QueryVideoById(int64(video))
			results = append(results, VideoParams(result))
		}
		log.Println("like list getting successfully!")
		return results
	}

}

/*获取用户userId喜欢的视频数量*/
func (like *LikeServiceImpl) LikeVideoCount(userId int64) (int64, error) {
	key_userId := utils.LikeUserKey + strconv.FormatInt(userId, 10)
	//先判断key_userId键值是否在缓存中
	if n, err := redis.RDb.Exists(redis.Ctx, key_userId).Result(); n > 0 { //key_userId键值在缓存中
		if err != nil {
			log.Println("redis Query failed")
			return 0, err
		} else { //查询成功
			result, err := redis.RDb.SCard(redis.Ctx, key_userId).Result() //获取key_userId键值有几个val
			if err != nil {
				log.Println("redis Query failed")
				return 0, err
			}
			//减去添加的默认值
			return result - 1, nil
		}
	} else { //key_userId键值不在缓存中，需要把MySQL中的数据添加到缓存中
		redis.RDb.SAdd(redis.Ctx, key_userId, utils.MyDefault).Result()
		redis.RDb.Expire(redis.Ctx, key_userId, utils.LikeUserKeyTTL).Result()

		//把数据库中的当前用户点赞的videoId全部添加到缓存中
		likevideoIdList, err1 := dao.GetLikeVideoIdList(userId)
		if err1 != nil {
			log.Println("Failed to get the likes video id list")
			redis.RDb.Del(redis.Ctx, key_userId)
			return 0, err1
		}
		for _, videoId := range likevideoIdList {
			strvideoId := strconv.FormatInt(videoId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RDb.SAdd(redis.Ctx, key_userId, strvideoId).Result(); err != nil {
				log.Println("Failed to add cache for videoId")
				redis.RDb.Del(redis.Ctx, key_userId)
				return 0, err
			}
		}

		return int64(len(likevideoIdList)), nil
	}

}

// func (like *LikeServiceImpl) IsLike(videoId int64, userId int64) bool {
// 	key_userId := utils.LikeUserKey + strconv.FormatInt(userId, 10)
// 	strVideoId := strconv.FormatInt(videoId, 10)
// 	isLike, _ := redis.RDb.SIsMember(redis.Ctx, key_userId, strVideoId).Result()
// 	log.Println("the statis return from redis is", isLike)
// 	return isLike

// if n, _ := redis.RDb.Exists(redis.Ctx, key_userId).Result(); n <= 0 { //如果key_userId存在缓存中
// 	log.Println("user not in redis")

// } else {
// 	log.Println("user in redis")
// 	videoIdList_database, _ := dao.GetLikeVideoIdList(userId)
// 	log.Println("the videoidlist find in database is", videoIdList_database)

// 	videoIdList, _ := redis.RDb.SMembers(redis.Ctx, key_userId).Result()
// 	log.Println("the videoidlist find in redis is", videoIdList)

// 	isLike, _ := redis.RDb.SIsMember(redis.Ctx, key_userId, strVideoId).Result()
// 	log.Println("the statis return from redis is", isLike)
// }
// return true
// }

// 	if err != nil {
// 		log.Println("Failed to get the likes video id list")
// 		return false
// 	}
// 	for _, vId := range videoIdList {
// 		if vId == videoId {
// 			return true
// 		}
// 	}
// 	return false
// }

/*判断用户userId是否点赞视频videoId*/
func (like *LikeServiceImpl) IsLike(videoId int64, userId int64) bool {
	strUserId := strconv.FormatInt(userId, 10)
	strVideoId := strconv.FormatInt(videoId, 10)
	key_userId := utils.LikeUserKey + strconv.FormatInt(userId, 10)
	key_videoId := utils.LikeVideokey + strconv.FormatInt(videoId, 10)

	if n, err := redis.RDb.Exists(redis.Ctx, key_userId).Result(); n > 0 { //如果key_userId存在缓存中
		if err != nil {
			log.Println("Redis query failed")
			return false
		}
		isLike, err := redis.RDb.SIsMember(redis.Ctx, key_userId, strVideoId).Result()
		if err != nil {
			log.Println("Redis query failed")
			return false
		}
		return isLike
	} else { //如果key_userId不存在缓存中，查询key_videoId是否在缓存中
		if n, err := redis.RDb.Exists(redis.Ctx, key_videoId).Result(); n > 0 { //如果key_userId存在缓存中
			if err != nil {
				log.Println("Redis query failed")
				return false
			}
			isLike, err := redis.RDb.SIsMember(redis.Ctx, key_videoId, strUserId).Result()
			if err != nil {
				log.Println("Redis query failed")
				return false
			}
			return isLike
		} else { //如果key_userId key_videoId都不存在缓存当中,就从数据库中读取然后添加到缓存中

			videoIdList, err := dao.GetLikeVideoIdList(userId)
			if err != nil {
				log.Println("Failed to get the likes video id list")
				return false
			}
			for _, vId := range videoIdList {
				if vId == videoId {
					//如果存在喜欢关系，那么添加缓存
					redis.RDb.SAdd(redis.Ctx, key_userId, utils.MyDefault).Result()
					redis.RDb.Expire(redis.Ctx, key_userId, utils.LikeUserKeyTTL).Result()
					//把数据库中的当前用户点赞的videoId全部添加到缓存中
					videoIdList, err1 := dao.GetLikeVideoIdList(userId)
					if err1 != nil {
						log.Println("Failed to get the likes video id list")
						redis.RDb.Del(redis.Ctx, key_userId)
						return false
					}
					for _, videoId := range videoIdList {
						strvideoId := strconv.FormatInt(videoId, 10)
						//如果出现一次不对的就把这个键值删除
						if _, err := redis.RDb.SAdd(redis.Ctx, key_userId, strvideoId).Result(); err != nil {
							log.Println("Failed to add cache for videoId")
							redis.RDb.Del(redis.Ctx, key_userId)
							return false
						}
					}
					isLike, err := redis.RDb.SIsMember(redis.Ctx, key_userId, strVideoId).Result()
					if err != nil {
						log.Println("Redis query failed")
						return false
					}

					redis.RDb.SAdd(redis.Ctx, key_videoId, utils.MyDefault).Result()
					redis.RDb.Expire(redis.Ctx, key_videoId, utils.LikeVideoKeyTTL).Result()

					//把数据库中给当前视频的点赞的userId全部添加到缓存中
					userIdList, err1 := dao.GetLikeUserIdList(videoId)
					if err1 != nil {
						log.Println("Failed to get video id like user list")
						redis.RDb.Del(redis.Ctx, key_videoId)
						return false
					}
					for _, userId := range userIdList {
						struserid := strconv.FormatInt(userId, 10)
						//如果出现一次不对的就把这个键值删除
						if _, err := redis.RDb.SAdd(redis.Ctx, key_videoId, struserid).Result(); err != nil {
							log.Println("Failed to add cache for videoId")
							redis.RDb.Del(redis.Ctx, key_videoId)
							return false
						}
					}

					return isLike

				}
			}

		}
	}

	return false
}

/*获取视频videoId的点赞数*/
func (like *LikeServiceImpl) CountLikes(videoId int64) int64 {
	key_videoId := utils.LikeVideokey + strconv.FormatInt(videoId, 10)
	//var result int64
	//result := 1
	//如果键值key_videoId在缓存中
	if n, err := redis.RDb.Exists(redis.Ctx, key_videoId).Result(); n > 0 {
		if err != nil {
			log.Println("Redis query failed")
			return -1
		}
		result, err1 := redis.RDb.SCard(redis.Ctx, key_videoId).Result()
		if err1 != nil {
			log.Println("Redis query failed")
			return -1
		}

		return result - 1
	} else { //如果键值key_videoId不在缓存中
		cnt, err := dao.CountLikes(videoId)
		if err != nil {
			log.Println("count from db error:", err)
			return 0
		}

		redis.RDb.SAdd(redis.Ctx, key_videoId, utils.MyDefault).Result()
		redis.RDb.Expire(redis.Ctx, key_videoId, utils.LikeVideoKeyTTL).Result()
		//把数据库中的点赞该视频的userid全部添加到缓存中
		userIdList, err1 := dao.GetLikeUserIdList(videoId)
		if err1 != nil {
			log.Println("Get likes video id list")
			redis.RDb.Del(redis.Ctx, key_videoId)
			return -1
		}
		for _, videoId := range userIdList {
			struserId := strconv.FormatInt(videoId, 10)
			//如果出现一次不对的就把这个键值删除
			if _, err := redis.RDb.SAdd(redis.Ctx, key_videoId, struserId).Result(); err != nil {
				log.Println("Failed to add cache for videoId")
				redis.RDb.Del(redis.Ctx, key_videoId)
				return -1
			}
		}
		return cnt
	}
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
