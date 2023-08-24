package service

import (
	"github.com/Shanwu404/TikTokLite/middleware/redis"
	"github.com/Shanwu404/TikTokLite/utils"
	"log"
	"strconv"

	"github.com/Shanwu404/TikTokLite/dao"
)

type CommentServiceImpl struct{}

func NewCommentService() CommentService {
	return &CommentServiceImpl{}
}

func (CommentServiceImpl) QueryCommentsByVideoId(id int64) []CommentParams {
	comments, err := dao.QueryCommentsByVideoId(id)
	if err != nil {
		log.Println("error:", err.Error())
	}
	results := make([]CommentParams, 0, len(comments))
	for i := range comments {
		results = append(results, CommentParams(comments[i]))
	}
	log.Println("Query comments successfully!")
	return results
}

func (CommentServiceImpl) PostComment(comment CommentParams) (int64, int32, string) {
	commentNew, err := dao.InsertComment(dao.Comment(comment))
	if err != nil {
		return -1, 1, "Post comment failed!"
	}
	if flag := CommentInsertRedis(comment.VideoId, comment.Id); !flag {
		log.Println("Insert redis failed!")
	}
	return commentNew.Id, 0, "Post comment successfully!"
}

func (CommentServiceImpl) DeleteComment(id int64) (int32, string) {
	flag := dao.DeleteComment(id)
	if flag == false {
		return 1, "Delete comment failed!"
	}
	return 0, "Delete comment successfully!"
}

func (CommentServiceImpl) CountComments(id int64) int64 {
	cnt, err := dao.CountComments(id)
	if err != nil {
		log.Println("count from db error:", err)
		return 0
	}
	log.Println("count comments successfully!")
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
