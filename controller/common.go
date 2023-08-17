package controller

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type UserInfo struct {
	Id              int64  `json:"id"`               // 用户ID
	Username        string `json:"name"`             // 用户名
	FollowCount     int64  `json:"follow_count"`     // 关注数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝数
	IsFollow        bool   `json:"is_follow"`        // 是否关注
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 背景图片
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  int64  `json:"total_favorited"`  // 获赞数
	WorkCount       int64  `json:"work_count"`       // 作品数
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
}

type CommentInfo struct {
	Id         int64
	User       UserInfo
	Content    string
	CreateDate string
}

type VideoInfo struct {
	User          UserInfo
	PlayURL       string
	CoverURL      string
	Title         string
	FavoriteCount int64
	CommentCount  int64
	IsFavorite    int64
}
