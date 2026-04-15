package middleware

import (
	"context"
	"log"
	"net/http"
	"restServer/response"
	"strings"
)

const usernameKey contextKey = "username"

// TokenValidator token 验证函数类型
type TokenValidator func(token string) (string, bool)

var tokenValidator TokenValidator

// SetTokenValidator 设置 token 验证函数
func SetTokenValidator(validator TokenValidator) {
	tokenValidator = validator
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//从 Header里面读取Authorization
		authHeader := r.Header.Get("Authorization")

		//检查是否提供了token
		if authHeader == "" {
			requestID := GetRequestID(r.Context())
			log.Printf("[Auth Failed] RequestID: %s - No token provided", requestID)
			response.Error(w, http.StatusUnauthorized, 40101, "没有提供认证令牌")
			return
		}

		//解析Bearer token
		//格式："Bearer <token>"
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			// 没有 "Bearer " 前缀
			requestID := GetRequestID(r.Context())
			log.Printf("[Auth Failed] RequestID:%v-Invalid token format", requestID)
			response.Error(w, http.StatusUnauthorized, 40102, "认证令牌格式错误")
			return
		}

		// 检查 token 是否为空
		if token == "" {
			requestID := GetRequestID(r.Context())
			log.Printf("[Auth Failed] RequestID: %s - Empty token", requestID)
			response.Error(w, http.StatusUnauthorized, 40102, "认证令牌不能为空")
			return
		}

		// 检查 tokenValidator 是否已设置
		if tokenValidator == nil {
			requestID := GetRequestID(r.Context())
			log.Printf("[Auth Error] RequestID: %s - Token validator not configured", requestID)
			response.Error(w, http.StatusInternalServerError, 50001, "认证服务未配置")
			return
		}

		//校验 token
		username, valid := tokenValidator(token)
		if !valid {
			requestID := GetRequestID(r.Context())
			log.Printf("[Auth Failed] RequestID: %s - Invalid token", requestID)
			response.Error(w, http.StatusUnauthorized, 40103, "认证令牌无效")
			return
		}

		//token 有效，将用户名存入 context
		requestID := GetRequestID(r.Context())
		log.Printf("[Auth Success] RequestID: %s - User: %s", requestID, username)

		ctx := context.WithValue(r.Context(), usernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUsername 从 context 中获取用户名
func GetUsername(ctx context.Context) string {
	if username, ok := ctx.Value(usernameKey).(string); ok {
		return username
	}
	return ""
}
