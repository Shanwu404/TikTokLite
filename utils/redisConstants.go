package utils

import (
	"math/rand"
	"time"
)

const Day = time.Hour * 24
const Month = Day * 30

var HourRandnum int64 = rand.Int63n(24)
var DayRandnum int64 = rand.Int63n(30)

// CommentCommentKey key:commentId value:videoId
const CommentCommentKey = "comment:comment:"

// CommentVideoKey key:videoId value:commentIds
const CommentVideoKey = "comment:video:"

// 喜欢模块
const LikeUserKey = "like:user:"

const LikeVideokey = "like:video:"

const MyDefault = -1

var LikeUserKeyTTL = Day*7 + Day*time.Duration(DayRandnum)

var LikeVideoKeyTTL = Day*7 + Day*time.Duration(DayRandnum)

/*------------------------ user模块-----------------------------*/

// UserIdKey key:userId value:User_struct
const UserIdKey = "user:id:"

var UserIdKeyTTL = Day*7 + Day*time.Duration(DayRandnum)

// UserNameKey key:username value:userId
const UserNameKey = "user:name:"

var UserNameKeyTTL = Day*7 + Day*time.Duration(DayRandnum)

/*------------------------ relation模块--------------------------*/

// RelationFollowKey key:userId value:followId
const RelationFollowKey = "relation:follow:"

var RelationFollowKeyTTL = Day*7 + Day*time.Duration(DayRandnum)

// RelationFollowCntKey key:userId value:followCnt
const RelationFollowCntKey = "relation:followCnt:"

var RelationFollowCntKeyTTL = Day*7 + Day*time.Duration(DayRandnum)

// RelationFollowerCntKey key:userId value:followerCnt
const RelationFollowerCntKey = "relation:followerCnt:"

var RelationFollowerCntKeyTTL = Day*7 + Day*time.Duration(DayRandnum)
