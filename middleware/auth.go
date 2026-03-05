package middleware

import (
	"log"
	"net/http"
	"restServer/response"
	"strings"
)

// 硬编码的有效 token（实际项目中应该从数据库或配置中读取）
const validToken = "my-secret-token-123456"

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

		//校验 token
		if token != validToken {
			requestID := GetRequestID(r.Context())
			log.Printf("[Auth Failed] RequestID: %s - Invalid token: %s", requestID, token)
			response.Error(w, http.StatusUnauthorized, 40103, "认证令牌无效")
			return
		}

		//token 有效，继续处理
		requestID := GetRequestID(r.Context())
		log.Printf("[Auth Success] RequestID: %s - Token validated", requestID)
		next.ServeHTTP(w, r)

	})
}
