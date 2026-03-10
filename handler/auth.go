package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"restServer/model"
	"restServer/response"
	"sync"
)

// 内存存储
var (
	users  = make(map[string]*model.User) // username -> User
	tokens = make(map[string]string)      // token -> username
	mu     sync.RWMutex
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

// Register 用户注册接口
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, 405, "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON 格式错误: "+err.Error())
		return
	}
	defer r.Body.Close()

	// 参数校验
	if req.Username == "" || req.Password == "" {
		response.BadRequest(w, "用户名和密码不能为空")
		return
	}
	if len(req.Username) < 3 || len(req.Username) > 20 {
		response.BadRequest(w, "用户名长度必须在3-20个字符之间")
		return
	}
	if len(req.Password) < 6 {
		response.BadRequest(w, "密码长度不能少于6位")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// 检查用户是否已存在
	if _, exists := users[req.Username]; exists {
		response.Error(w, http.StatusConflict, 40901, "用户名已存在")
		return
	}

	// 存储用户
	users[req.Username] = &model.User{
		Name:     req.Username,
		Password: req.Password,
	}

	response.Success(w, map[string]any{
		"message":  "注册成功",
		"username": req.Username,
	})
}

// Login 用户登录接口
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, 405, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON 格式错误: "+err.Error())
		return
	}
	defer r.Body.Close()

	// 参数校验
	if req.Username == "" || req.Password == "" {
		response.BadRequest(w, "用户名和密码不能为空")
		return
	}

	mu.RLock()
	user, exists := users[req.Username]
	mu.RUnlock()

	// 验证用户名和密码
	if !exists || user.Password != req.Password {
		response.Error(w, http.StatusUnauthorized, 40101, "用户名或密码错误")
		return
	}

	// 生成 token
	token := generateToken()

	mu.Lock()
	tokens[token] = req.Username
	mu.Unlock()

	response.Success(w, LoginResponse{
		Token:    token,
		Username: req.Username,
	})
}

// generateToken 生成随机 token
func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GetUsernameByToken 根据 token 获取用户名
func GetUsernameByToken(token string) (string, bool) {
	mu.RLock()
	defer mu.RUnlock()
	username, exists := tokens[token]
	return username, exists
}
