package service

import (
	"errors"
	"log"
	"unicode"

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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Username does not exist: ", username)
			return dao.User{}, err
		}
		log.Println("Error querying user by name:", username, err)
		return dao.User{}, err
	}
	user.Password = "" // 屏蔽密码
	log.Println("Query user successfully! User queried by name:", user)
	return *user, nil
}

// QueryUserByID 根据id获取User对象 屏蔽密码
func (us *UserServiceImpl) QueryUserByID(id int64) (dao.User, error) {
	log.Println("Querying user by ID:", id)
	user, err := dao.QueryUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("User ID not found:", id)
			return dao.User{}, err
		}
		log.Println("Error querying user by ID:", id, "-", err.Error())
		return dao.User{}, err
	}
	user.Password = "" // 屏蔽密码
	log.Println("Query user successfully! User queried by ID:", user)
	return *user, nil
}

// Register 用户注册，返回注册用户ID，状态码和状态信息
func (us *UserServiceImpl) Register(username string, password string) (int64, int32, string) {

	// 验证用户名和密码的合法性
	if !isValidUsername(username) {
		log.Println("Invalid username format:", username)
		return -1, 1, "Invalid username format!"
	}
	if !isValidPassword(password) {
		log.Println("Invalid password format")
		return -1, 1, "Invalid password format!"
	}

	log.Println("Registering user:", username)
	user, _ := dao.QueryUserByUsername(username)

	if user != nil {
		return -1, 1, "User already exist!"
	}

	encoderPassword, err := HashEncode(password)
	if err != nil {
		log.Println("Password encoding error:", err)
		return -1, 1, "Incorrect password format!"
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

	// 验证用户名和密码的合法性
	if !isValidUsername(username) {
		log.Println("Invalid username format:", username)
		return 1, "Invalid username format!"
	}
	if !isValidPassword(password) {
		log.Println("Invalid password format")
		return 1, "Invalid password format!"
	}

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

// IsUserIdExist 查询用户ID是否存在
func (us *UserServiceImpl) IsUserIdExist(id int64) bool {
	log.Println("Checking if user ID exists:", id)
	isExisted := dao.IsUserIdExist(id)
	log.Printf("User ID %d exists: %t\n", id, isExisted)
	return isExisted
}

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
		log.Println("Password comparison error:", err)
		return false
	}
	return true
}
