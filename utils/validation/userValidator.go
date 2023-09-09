package validation

import "unicode"

func IsValidUsername(username string) bool {
	// 用户名长度限制为3-12个字符
	const minUsernameLength = 3
	const maxUsernameLength = 32
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

func IsValidPassword(password string) bool {
	// 密码长度限制为3-12个字符
	const minPasswordLength = 5
	const maxPasswordLength = 32
	length := len(password)

	if length < minPasswordLength || length > maxPasswordLength {
		return false
	}

	// 密码只包括 ASCII 字母、数字和标点符号
	for _, ch := range password {
		if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') && !unicode.IsPunct(ch) {
			return false
		}
	}

	return true
}