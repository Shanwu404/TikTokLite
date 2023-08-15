package controller

import (
	"github.com/Shanwu404/TikTokLite/dao"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type User struct {
	dao.UserResp
}

type CommentInfo struct {
	Id         int64
	User       dao.UserResp
	Content    string
	CreateDate string
}
