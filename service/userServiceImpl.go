package service

import (
	"errors"
	"log"

	"github.com/Shanwu404/TikTokLite/dao"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
}

func NewUserService() UserService {
	return &UserServiceImpl{}
}

// QueryUserByName 根据name获取User对象
func (us *UserServiceImpl) QueryUserByUsername(username string) (dao.User, error) {
	log.Println("Querying user by name:", username)
	user, err := dao.QueryUserByUsername(username)
	if err != nil {
		log.Println("Error querying user by name:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(username, " not found!")
		}
		return dao.User{}, err
	}
	user.Password = "" // 屏蔽密码
	log.Println("Query user successfully! User queried by name:", user)
	return *user, nil
}

// QueryUserByID 根据id获取User对象 屏蔽密码
func (us *UserServiceImpl) QueryUserByID(id uint64) (dao.User, error) {
	log.Println("Query user successfully! Querying user by ID:", id)
	user, err := dao.QueryUserByID(id)
	if err != nil {
		log.Println("Error querying user by ID:", id, err)
		return dao.User{}, err
	}
	user.Password = "" // 屏蔽密码
	log.Println("User queried by ID:", user)
	return *user, nil
}

func (us *UserServiceImpl) QueryUserRespByID(id uint64) (dao.UserResp, error) {
	userInfo := dao.UserResp{}
	user, err := us.QueryUserByID(id)
	if err != nil {
		log.Println(err.Error())
		return userInfo, err
	}
	log.Println("User queried by ID:", user)
	userInfo.Id = user.ID             // 用户ID
	userInfo.Username = user.Username // 用户名
	userInfo.FollowerCount = 0        // 粉丝数接口待实现
	userInfo.FollowCount = 0          // 关注数接口待实现
	userInfo.IsFollow = false         // 是否关注接口待实现
	userInfo.Avatar = ""              // 用户头像接口待实现
	userInfo.BackgroundImage = ""     // 背景图片接口待实现
	userInfo.Signature = ""           // 个人简介接口待实现
	userInfo.TotalFavorited = 0       // 获赞数接口待实现
	userInfo.FavoriteCount = 0        // 喜欢数接口待实现
	userInfo.WorkCount = 0            // 作品数接口待实现
	return userInfo, nil
}

// Register 用户注册，返回注册用户ID，状态码和状态信息
func (us *UserServiceImpl) Register(username string, password string) (uint64, int32, string) {
	log.Println("Registering user:", username)
	user, _ := dao.QueryUserByUsername(username)

	// 先user返回值不为空指针，再去判断是否username已存在
	if user != nil && username == user.Username {
		return 0, 1, "User already exist!"
	}

	encoderPassword, err := HashEncode(password)
	if err != nil {
		log.Println("Password encoding error:", err)
		return 0, 1, "Incorrect password format!"
	}
	newUser := &dao.User{ // 创建一个指向User的指针
		Username: username,
		Password: encoderPassword,
	}
	err = dao.InsertUser(newUser)
	if err != nil {
		log.Println("User registration error:", err)
		return 0, 1, "User registration failed!"
	}

	log.Println("User registered successfully:", newUser)
	return newUser.ID, 0, "Register successfully!"
}

// Login 用户登录，返回状态码和状态信息
func (us *UserServiceImpl) Login(username string, password string) (int32, string) {
	log.Println("Attempting login for user:", username)
	// 注意要使用dao.QueryUserByUsername接口查询，否则密码为空
	user, err := dao.QueryUserByUsername(username)
	if err != nil {
		log.Println("User login error:", err)
		return 1, "User doesn't exist!"
	}
	arePasswordsEqual := ComparePasswords(user.Password, password)
	if arePasswordsEqual {
		return 0, "Login success"
	} else {
		return 1, "Username or Password error"
	}
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
		log.Println("Password comparison error:", err)
		return false
	}
	return true
}