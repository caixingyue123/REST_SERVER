package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter 包装 http.ResponseWriter 以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Logger 日志中间件 - 记录请求的方法、路径、耗时
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		//创建自定义ResponseWriter来捕获状态码
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		//调用下一个处理器
		next.ServeHTTP(rw, r)

		//记录日志
		duration := time.Since(start)
		log.Printf("[%s] %s %s - Status: %d - Duration: %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			rw.statusCode,
			duration,
		)
	})
}

func (rw *responseWriter) WriteHeader(code int){
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}