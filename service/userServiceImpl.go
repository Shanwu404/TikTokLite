package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Shanwu404/TikTokLite/dao"
	"github.com/Shanwu404/TikTokLite/log/logger"
	"github.com/Shanwu404/TikTokLite/middleware/redis"
	"github.com/Shanwu404/TikTokLite/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	relationService RelationService
	videoService    VideoService
	likeService     LikeService
}

func NewUserService() UserService {
	return &UserServiceImpl{
		relationService: NewRelationService(),
		videoService:    NewVideoService(),
		likeService:     NewLikeService(),
	}
}

// QueryUserByUsername 根据name获取User对象
func (us *UserServiceImpl) QueryUserByUsername(username string) (dao.User, error) {
	logger.Infoln("Querying user by name: ", username)

	// 尝试从Redis中获取用户ID
	redisNameKey := utils.UserNameKey + username
	userIDStr, err := redis.RDb.Get(redis.Ctx, redisNameKey).Result()
	if err == nil && userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			logger.Errorln("Error parsing user ID from Redis:", err)
		} else {
			return us.QueryUserByID(userID)
		}
	}

	// 如果Redis中没有用户ID，则从数据库中获取
	user, err := dao.QueryUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorln("Username does not exist:", username)
			return dao.User{}, err
		}
		logger.Errorln("Error querying user by name:", username, "-", err.Error())
		return dao.User{}, err
	}
	user.Password = "" // 屏蔽密码

	// 将用户ID存入Redis
	redis.RDb.Set(redis.Ctx, redisNameKey, user.ID, utils.UserNameKeyTTL)

	// 将用户信息存入Redis
	redisIdKey := utils.UserIdKey + strconv.FormatInt(user.ID, 10)
	userBytes, err := json.Marshal(user)
	if err == nil {
		redis.RDb.Set(redis.Ctx, redisIdKey, userBytes, utils.UserIdKeyTTL)
	}
	logger.Infoln("Query user successfully (MySQL)! User queried by name: ", user.Username)
	return user, nil
}

// QueryUserByID 根据id获取User对象 屏蔽密码
func (us *UserServiceImpl) QueryUserByID(id int64) (dao.User, error) {
	logger.Infoln("Querying user by ID:", id)

	// 尝试从Redis中获取用户信息
	redisIdKey := utils.UserIdKey + strconv.FormatInt(id, 10)
	userData, err := redis.RDb.Get(redis.Ctx, redisIdKey).Result()
	if err == nil && userData != "" {
		// 解析Redis中的数据到User对象
		var user dao.User
		err := json.Unmarshal([]byte(userData), &user)
		if err == nil {
			logger.Infoln("Query user successfully (Redis)! User queried by ID: ", id)
			return user, nil
		}
	}

	// 如果Redis中没有用户信息或解析失败，则从数据库中获取
	user, err := dao.QueryUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorln("User ID not found:", id)
			return dao.User{}, err
		}
		logger.Errorln("Error querying user by ID:", id, "-", err.Error())
		return dao.User{}, err
	}
	user.Password = "" // 屏蔽密码

	// 将用户信息存入Redis
	userBytes, err := json.Marshal(user)
	if err == nil {
		redis.RDb.Set(redis.Ctx, redisIdKey, userBytes, utils.UserIdKeyTTL)
	}
	logger.Infoln("Query user successfully (MySQL)! User queried by ID: ", user)
	return user, nil
}

// Register 用户注册，返回注册用户ID，状态码和状态信息
func (us *UserServiceImpl) Register(username string, password string) (int64, int32, string) {
	logger.Infoln("Registering user:", username)

	// 获取分布式锁
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	value := strconv.Itoa(r.Int())
	lock := redis.Lock("register_lock:"+username, value)
	if !lock {
		logger.Errorln("Failed to get distributed lock")
		return -1, 1, "Username registration underway. Please try later."
	}
	defer redis.Unlock("register_lock:" + username)

	// 检查用户名是否已存在
	user, err := us.QueryUserByUsername(username)
	if err == nil && user.Username != "" {
		logger.Errorln("User already exists:", username)
		return -1, 1, "User already exist!"
	}

	encoderPassword, err := HashEncode(password)
	if err != nil {
		logger.Errorln("Password encoding error:", err)
		return -1, 1, "Incorrect password format!"
	}

	newUser := dao.User{
		Username: username,
		Password: encoderPassword,
	}

	newUser.ID, err = dao.InsertUser(newUser)
	if err != nil {
		logger.Errorln("User registration error:", err)
		return 0, 1, "User registration failed!"
	}

	// 将用户信息存入Redis
	redisIdKey := "user:id:" + strconv.FormatInt(newUser.ID, 10)
	redisNameKey := "user:name:" + newUser.Username
	userBytes, err := json.Marshal(newUser)
	if err == nil {
		redis.RDb.Set(redis.Ctx, redisIdKey, userBytes, utils.UserIdKeyTTL)
		redis.RDb.Set(redis.Ctx, redisNameKey, newUser.ID, utils.UserNameKeyTTL)
	}

	logger.Infoln("User registered successfully:", newUser.Username)
	return newUser.ID, 0, "Register successfully!"
}

// Login 用户登录，返回用户ID状态码和状态信息
func (us *UserServiceImpl) Login(username string, password string) (int64, int32, string) {
	logger.Infoln("Attempting login for user:", username)

	user, err := dao.QueryUserByUsername(username)
	if err != nil {
		logger.Errorln("User login error:", err)
		return -1, 1, "User doesn't exist!"
	}

	arePasswordsEqual := ComparePasswords(user.Password, password)
	if arePasswordsEqual {
		return user.ID, 0, "Login success"
	} else {
		return -1, 1, "Username or Password error"
	}
}

// IsUserIdExist 查询用户ID是否存在
func (us *UserServiceImpl) IsUserIdExist(id int64) bool {
	logger.Infoln("Checking if user ID exists:", id)

	// 尝试从Redis中获取用户信息
	redisIdKey := "user:id:" + strconv.FormatInt(id, 10)
	userData, err := redis.RDb.Get(redis.Ctx, redisIdKey).Result()
	if err == nil && userData != "" {
		// 说明Redis中存在该用户信息
		logger.Infof("User ID %d exists (Redis)\n", id)
		return true
	}

	// 如果Redis中没有用户信息，则从数据库中获取
	user, err := dao.QueryUserByID(id)
	if err != nil {
		logger.Errorln("User ID not found:", id, "-", err.Error())
		return false
	}
	user.Password = "" // 屏蔽密码

	// 将用户信息存入Redis
	userBytes, err := json.Marshal(user)
	if err == nil {
		redis.RDb.Set(redis.Ctx, redisIdKey, userBytes, utils.UserIdKeyTTL)
	}
	logger.Infof("User ID %d exists (MySQL)\n", id)
	return true
}

// QueryUserInfoByID 根据用户ID查询用户信息
func (us *UserServiceImpl) QueryUserInfoByID(userId int64) (UserInfoParams, error) {
	logger.Infoln("Querying userinfo by ID:", userId)

	// 判断用户ID是否存在
	if isExisted := us.IsUserIdExist(userId); !isExisted {
		return UserInfoParams{}, fmt.Errorf("user doesn't exist")
	}

	user, _ := us.QueryUserByID(userId)
	followCount, _ := us.relationService.CountFollows(userId)
	followerCount, _ := us.relationService.CountFollowers(userId)
	favoriteCount, _ := us.likeService.LikeVideoCount(userId)
	totalFavorited := us.likeService.TotalFavorited(userId)
	videos := us.videoService.GetVideoListByUserId(userId)
	workCount := int64(len(videos))
	// workCount := int64(10)

	userInfo := UserInfoParams{
		Id:              user.ID,
		Username:        user.Username,
		FollowCount:     followCount,
		FollowerCount:   followerCount,
		IsFollow:        false, // 注意这个值需要根据具体情况修改
		Avatar:          "https://mary-aliyun-img.oss-cn-beijing.aliyuncs.com/typora/202308171029672.jpg",
		BackgroundImage: "https://mary-aliyun-img.oss-cn-beijing.aliyuncs.com/typora/202308171007006.jpg",
		Signature:       "这个人很懒，什么都没有留下",
		TotalFavorited:  totalFavorited,
		WorkCount:       workCount,
		FavoriteCount:   favoriteCount,
	}
	return userInfo, nil
}

/*------------------------ 以下为工具函数 ------------------------*/
// HashEncode 加密密码
func HashEncode(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePasswords 验证密码，password1为加密的密码，password2为待验证的密码
func ComparePasswords(password1 string, password2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
	if err != nil {
		logger.Errorln("Password comparison error:", err)
		return false
	}
	return true
}
