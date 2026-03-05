package model

import (
	"fmt"
	"regexp"
)

// User 用户模型
type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
}

// Validate 参数校验
func (u *User) Validate() error {
	// 校验用户名
	if u.Name == "" {
		return fmt.Errorf("用户的姓名不能为空")
	}
	if len(u.Name) < 3 || len(u.Name) > 20 {
		return fmt.Errorf("用户名的长度只能在3-20个字符之间")
	}

	// 校验密码
	if len(u.Password) < 8 {
		return fmt.Errorf("密码长度不能少于8位")
	}

	matched, _ := regexp.MatchString("^[a-zA-Z0-9]+$", u.Password)
	if !matched {
		return fmt.Errorf("密码只能包含数字和大小写字母")
	}

	// 校验邮箱（如果提供）
	if u.Email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(u.Email) {
			return fmt.Errorf("邮箱格式不正确")
		}
	}

	return nil
}
