package service

import (
	"log"
	"strconv"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/log/logger"
	"github.com/Shanwu404/TikTokLite/middleware/rabbitmq"
	"github.com/Shanwu404/TikTokLite/middleware/redis"
	"github.com/Shanwu404/TikTokLite/utils"
)

type CommentServiceImpl struct{}

func NewCommentService() CommentService {
	return &CommentServiceImpl{}
}

func (CommentServiceImpl) QueryCommentsByVideoId(id int64) []CommentParams {
	comments, err := dao.QueryCommentsByVideoId(id)
	if err != nil {
		logger.Errorln("error:", err.Error())
	}
	results := make([]CommentParams, 0, len(comments))
	for i := range comments {
		results = append(results, CommentParams(comments[i]))
	}
	logger.Infoln("Query comments successfully!")
	return results
}

func (CommentServiceImpl) PostComment(comment CommentParams) (int64, int32, string) {
	commentNew, err := dao.InsertComment(dao.Comment(comment))
	if err != nil {
		return -1, 1, "Post comment failed!"
	}
	if flag := CommentInsertRedis(commentNew.VideoId, commentNew.Id); !flag {
		log.Println("Insert redis failed!")
	}
	return commentNew.Id, 0, "Post comment successfully!"
}

func (CommentServiceImpl) DeleteComment(id int64) (int32, string) {
	redisCommentKey := utils.CommentCommentKey + strconv.FormatInt(id, 10)
	if err := redis.RDb.Exists(redis.Ctx, redisCommentKey).Err(); err != nil {
		log.Println(err.Error())
	} else {
		videoId, _ := redis.RDb.Get(redis.Ctx, redisCommentKey).Result()
		redisVideoKey := utils.CommentVideoKey + videoId
		CommentDeleteRedis(redisCommentKey, redisVideoKey, id)
		rabbitmq.CommentDel.Producer(strconv.FormatInt(id, 10))
	}
	flag := dao.DeleteComment(id)
	if !flag {
		return 1, "Delete comment failed!"
	}
	return 0, "Delete comment successfully!"
}

func (CommentServiceImpl) CountComments(id int64) int64 {
	redisVideoKey := utils.CommentVideoKey + strconv.FormatInt(id, 10)
	cnt, err := redis.RDb.SCard(redis.Ctx, redisVideoKey).Result()
	if err != nil {
		log.Println("count from redis error:", err)
	}
	redis.RDb.Expire(redis.Ctx, redisVideoKey, redis.RandomTime())
	if cnt > 0 {
		return cnt
	}
	cnt, err = dao.CountComments(id)
	if err != nil {
		logger.Errorln("count from db error:", err)
		return 0
	}
	logger.Infoln("count comments successfully!")
	return cnt
}

func CommentInsertRedis(videoId int64, commentId int64) bool {
	redisVideoKey := utils.CommentVideoKey + strconv.FormatInt(videoId, 10)
	if _, err := redis.RDb.SAdd(redis.Ctx, redisVideoKey, commentId).Result(); err != nil {
		log.Println("Insert key:video_id value:comment_id into redis failed!")
		redis.RDb.Del(redis.Ctx, redisVideoKey)
		return false
	}
	redis.RDb.Expire(redis.Ctx, redisVideoKey, redis.RandomTime()) // 更新过期时间
	redisCommentKey := utils.CommentCommentKey + strconv.FormatInt(commentId, 10)
	if _, err := redis.RDb.Set(redis.Ctx, redisCommentKey, videoId, redis.RandomTime()).Result(); err != nil {
		log.Println("Insert key:comment_id value:video_id into redis failed!")
		return false
	}
	log.Println("Insert record into redis successfully!")
	return true
}

func CommentDeleteRedis(key1 string, key2 string, id int64) bool {
	if err := redis.RDb.Del(redis.Ctx, key1).Err(); err != nil {
		log.Println(err.Error())
		return false
	}
	if err := redis.RDb.SRem(redis.Ctx, key2, id).Err(); err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
