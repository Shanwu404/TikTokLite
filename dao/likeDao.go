package dao

import (
	"github.com/Shanwu404/TikTokLite/log/logger"
)

type Like struct {
	ID      int64
	UserId  int64
	VideoId int64
}

// TableName 获取点赞表名
func (Like) TableName() string {
	return "likes"
}

// InsertLike 插入点赞数据
func InsertLike(likeData *Like) error {
	err := db.Model(&Like{}).Create(&likeData).Error
	if err != nil {
		logger.Errorln(err)
		return err
	}
	return nil
}

// DeleteLike 删除点赞数据
func DeleteLike(userId int64, videoId int64) error {
	err := db.Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).Delete(&Like{}).Error
	if err != nil {
		logger.Errorln(err)
		return err
	}
	return nil
}

// GetLikeVideoIdList 根据userId查询其点赞全部videoId
func GetLikeVideoIdList(userId int64) ([]int64, error) {
	var likeVideoIdList []int64
	err := db.Model(&Like{}).Where(map[string]interface{}{"user_id": userId}).Pluck("video_id", &likeVideoIdList).Error
	if err != nil {
		if err.Error() == "record not found" {
			logger.Errorln("there are no likeVideoIds")
			return likeVideoIdList, nil
		} else {
			logger.Errorln("get likeVideoIdList failed: ", err)
			return likeVideoIdList, err
		}
	}
	return likeVideoIdList, nil
}

// GetLikeUserIdList 根据videoId查询点赞该视频的全部user_id
func GetLikeUserIdList(videoId int64) ([]int64, error) {
	var likeUserIdList []int64
	err := db.Model(Like{}).Where(map[string]interface{}{"video_id": videoId}).Pluck("user_id", &likeUserIdList).Error
	if err != nil {
		logger.Errorln("get likeUserIdList failed: ", err)
		return nil, err
	} else {
		return likeUserIdList, nil
	}
}

func CountLikes(videoId int64) (int64, error) {
	var count int64
	err := db.Model(&Like{}).Where("video_id = ?", videoId).Count(&count).Error
	if err != nil {
		logger.Errorln(err)
		return -1, err
	}
	return count, nil
}
