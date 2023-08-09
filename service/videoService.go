package service

import (
	"time"

	"github.com/Shanwu404/TikTokLite/dao"
)

type VideoDetail struct {
	Id            uint64
	Author        AuthInfo
	PlayUrl       string
	CoverUrl      string
	FavoriteCount uint64
	CommentCount  uint64
	IsFavorite    bool
	Title         string
}

// 后面要移入用户模块
type UserInfo struct {
	ID              uint64
	Name            string
	FollowCount     uint64
	FollowerCount   uint64
	IsFollow        bool
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  string
	WorkCount       uint64
	FavoriteCount   uint64
}

type AuthInfo = UserInfo

type VideoService interface {
	// QueryVideoById 根据视频id获取视频
	QueryVideoById(videoID uint64) dao.Video

	// 根据视频id和查询用户id查询视频信息
	QueryVideoInfoByVideoId(videoId uint64, queryUserId uint64) (VideoDetail, time.Time)

	// GetVideoIdListByUserId 根据作者id查询视频id列表
	GetVideoIdListByUserId(authorId uint64) []uint64

	// 根据时间获取视频id列表
	GetMultiVideo(latestTime time.Time, count int) []dao.Video

	// InsertVideosTable 将video插入videos表内
	InsertVideosTable(video dao.Video) bool

	// 存储视频文件
	StoreVideo(data []byte, username string, fileName string) error

	// CountWorks 统计用户id的作品数
	CountWorks(id int64) uint
}
