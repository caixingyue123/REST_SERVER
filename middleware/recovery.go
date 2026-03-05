package middleware

import (
	"fmt"
	"log"
	"net/http"
	"restServer/response"
	"runtime/debug"
)

// Recovery 中间件 - 捕获 panic，避免服务崩溃
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				//获取request_id（如果有）
				requestID := GetRequestID(r.Context())

				//记录错误堆栈
				log.Printf("[PANIC RECOVERED] RequestID: %s, Error: %v\n%s",
					requestID,
					err,
					debug.Stack(),
				)
				//
				response.InternalError(w, fmt.Sprintf("服务器内部错误(RequestID:%v)", requestID))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
