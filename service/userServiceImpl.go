package service

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"
	"unicode"

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
	logger.Infoln("INFO: Querying user by name: ", username)

	// 尝试从Redis中获取用户ID
	redisNameKey := utils.UserNameKey + username
	userIDStr, err := redis.RDb.Get(redis.Ctx, redisNameKey).Result()
	if err == nil && userIDStr != "" {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Println("ERROR: Error parsing user ID from Redis:", err)
		} else {
			return us.QueryUserByID(userID)
		}
	}

	// 如果Redis中没有用户ID，则从数据库中获取
	user, err := dao.QueryUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Errorln("WARN: Username does not exist:", username)
			return dao.User{}, err
		}
		logger.Errorln("ERROR: Error querying user by name:", username, "-", err.Error())
		return dao.User{}, err
	}
	user.Password = "" // 屏蔽密码

	// 将用户ID存入Redis
	redis.RDb.Set(redis.Ctx, redisNameKey, user.ID, 24*time.Hour)

	// 将用户信息存入Redis
	redisDataKey := "user:id:" + strconv.FormatInt(user.ID, 10)
	userBytes, err := json.Marshal(user)
	if err == nil {
		redis.RDb.Set(redis.Ctx, redisDataKey, userBytes, 2*time.Hour)
	}
	logger.Infoln("INFO: Query user successfully (MySQL)! User queried by name: ", user.Username)
	return user, nil
}

// QueryUserByID 根据id获取User对象 屏蔽密码
func (us *UserServiceImpl) QueryUserByID(id int64) (dao.User, error) {
	logger.Infoln("INFO: Querying user by ID:", id)

	// 尝试从Redis中获取用户信息
	redisIdKey := utils.UserIdKey + strconv.FormatInt(id, 10)
	userData, err := redis.RDb.Get(redis.Ctx, redisIdKey).Result()
	if err == nil && userData != "" {
		// 解析Redis中的数据到User对象
		var user dao.User
		err := json.Unmarshal([]byte(userData), &user)
		if err == nil {
			log.Println("INFO: Query user successfully (Redis)! User queried by ID: ", id)
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
		redis.RDb.Set(redis.Ctx, redisIdKey, userBytes, 24*time.Hour)
	}
	logger.Infoln("INFO: Query user successfully (MySQL)! User queried by ID: ", user)
	return user, nil
}

// Register 用户注册，返回注册用户ID，状态码和状态信息
func (us *UserServiceImpl) Register(username string, password string) (int64, int32, string) {

	// 验证用户名和密码的合法性
	if !isValidUsername(username) {
		logger.Errorln("WARN: Invalid username format:", username)
		return -1, 1, "Invalid username format!"
	}
	if !isValidPassword(password) {
		logger.Errorln("WARN: Invalid password format")
		return -1, 1, "Invalid password format!"
	}

	logger.Infoln("INFO: Registering user:", username)

	// 获取分布式锁
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	value := strconv.Itoa(r.Int())
	lock := redis.Lock("register_lock:"+username, value)
	if !lock {
		log.Println("WARN: Failed to get distributed lock")
		return -1, 1, "Username registration underway. Please try later."
	}
	defer redis.Unlock("register_lock:" + username)

	// 检查用户名是否已存在
	user, err := us.QueryUserByUsername(username)
	if err == nil && user.Username != "" {
		log.Println("User already exists:", username)
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

	err = dao.InsertUser(newUser)
	if err != nil {
		logger.Errorln("ERROR: User registration error:", err)
		return 0, 1, "User registration failed!"
	}

	// 将用户信息存入Redis
	redisIdKey := "user:id:" + strconv.FormatInt(newUser.ID, 10)
	redisNameKey := "user:name:" + newUser.Username
	userBytes, err := json.Marshal(newUser)
	if err == nil {
		redis.RDb.Set(redis.Ctx, redisIdKey, userBytes, 24*time.Hour)
		redis.RDb.Set(redis.Ctx, redisNameKey, newUser.ID, 24*time.Hour)
	}

	logger.Infoln("INFO: User registered successfully:", newUser.Username)
	return newUser.ID, 0, "Register successfully!"
}

// Login 用户登录，返回状态码和状态信息
func (us *UserServiceImpl) Login(username string, password string) (int32, string) {
	logger.Infoln("INFO: Attempting login for user:", username)

	// 验证用户名和密码的合法性
	if !isValidUsername(username) {
		logger.Errorln("WARN: Invalid username format:", username)
		return 1, "Invalid username format!"
	}
	if !isValidPassword(password) {
		logger.Errorln("WARN: Invalid password format")
		return 1, "Invalid password format!"
	}

	user, err := dao.QueryUserByUsername(username)
	if err != nil {
		logger.Errorln("ERROR: User login error:", err)
		return 1, "User doesn't exist!"
	}

	arePasswordsEqual := ComparePasswords(user.Password, password)
	if arePasswordsEqual {
		return 0, "Login success"
	} else {
		return 1, "Username or Password error"
	}
}

// IsUserIdExist 查询用户ID是否存在
func (us *UserServiceImpl) IsUserIdExist(id int64) bool {
	log.Println("INFO: Checking if user ID exists:", id)

	// 尝试从Redis中获取用户信息
	redisIdKey := "user:id:" + strconv.FormatInt(id, 10)
	userData, err := redis.RDb.Get(redis.Ctx, redisIdKey).Result()
	if err == nil && userData != "" {
		// 说明Redis中存在该用户信息
		log.Printf("INFO: User ID %d exists (Redis)\n", id)
		return true
	}

	// 如果Redis中没有用户信息，则从数据库中获取
	user, err := dao.QueryUserByID(id)
	if err != nil {
		log.Println("WARN: User ID not found:", id)
		return false
	}
	user.Password = "" // 屏蔽密码

	// 将用户信息存入Redis
	userBytes, err := json.Marshal(user)
	if err == nil {
		redis.RDb.Set(redis.Ctx, redisIdKey, userBytes, 24*time.Hour)
	}
	log.Printf("INFO: User ID %d exists (MySQL)\n", id)
	return true
}

// QueryUserInfoByID 根据用户ID查询用户信息
func (us *UserServiceImpl) QueryUserInfoByID(userId int64) (UserInfoParams, error) {
	logger.Infoln("Querying userinfo by ID:", userId)

	user, _ := us.QueryUserByID(userId)
	followCount, _ := us.relationService.CountFollows(userId)
	followerCount, _ := us.relationService.CountFollowers(userId)
	favoriteCount, _ := us.likeService.LikeVideoCount(userId)
	totalFavorited := us.likeService.TotalFavorited(userId)
	// videos := us.videoService.GetVideoListByUserId(userId)
	// workCount := int64(len(videos))
	workCount := int64(10) // 临时设置

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

func isValidUsername(username string) bool {
	// 用户名长度限制为3-12个字符
	const minUsernameLength = 3
	const maxUsernameLength = 12
	length := len(username)

	// 检查长度是否在范围内
	if length < minUsernameLength || length > maxUsernameLength {
		return false
	}

	// 检查用户名是否只包含字母和数字
	for _, ch := range username {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			return false
		}
	}

	return true
}

func isValidPassword(password string) bool {
	// 密码长度限制为3-12个字符
	const minPasswordLength = 3
	const maxPasswordLength = 12
	length := len(password)

	if length < minPasswordLength || length > maxPasswordLength {
		return false
	}

	// 密码只包括字母、数字和标点符号
	for _, ch := range password {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && !unicode.IsPunct(ch) {
			return false
		}
	}

	return true
}

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
