// 数据校验
package controller

import (
	"time"

	"github.com/Shanwu404/TikTokLite/dao"
)

type feedReq struct {
	latest_time time.Time
}

type publishReq struct {
}

type authorInfo = dao.UserResp // 结构不合理

type videoDetail struct {
	Id            uint64     `json:"id"`
	Author        authorInfo `json:"author"`
	PlayUrl       string     `json:"play_url"`
	CoverUrl      string     `json:"cover_url"`
	FavoriteCount uint64     `json:"favorite_count"`
	CommentCount  uint64     `json:"comment_count"`
	IsFavorite    bool       `json:"is_favorite"`
	Title         string     `json:"title"`
}

type feedResp struct {
	Response
	NextTime  int64       `json:"next_time"`
	VideoList videoDetail `json:"video_list"`
}
