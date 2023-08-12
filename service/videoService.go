package service

import "time"

type VideoParams struct {
	ID          int64
	AuthorID    int64
	PlayURL     string
	CoverURL    string
	PublishTime time.Time
	Title       string
}

// 后面要移入用户模块
type UserParams struct {
	ID              int64
	Name            string
	FollowCount     int64
	FollowerCount   int64
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  string
	WorkCount       int64
	FavoriteCount   int64
}

type VideoService interface {
	// QueryVideoById 根据视频id获取视频
	QueryVideoById(videoID int64) VideoParams

	// 根据视频id和查询用户id查询视频信息
	// QueryVideoInfoByVideoId(videoId int64, queryUserId int64) (VideoParams, time.Time)

	// GetVideoIdListByUserId 根据作者id查询视频列表
	GetVideoListByUserId(authorId int64) []VideoParams

	// 根据时间获取视频id列表
	GetMultiVideoBefore(latestTime int64, count int) []VideoParams

	// InsertVideosTable 将video插入videos表内
	InsertVideosTable(video *VideoParams) error

	// 存储视频文件
	StoreVideo(data []byte, username string, fileName string) error

	// 统计用户id的作品数
	// WorkCount(id int64) int
}
