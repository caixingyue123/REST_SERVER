package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"
)

// contextKey 用于在 context 中存储值的类型
type contextKey string

const RequestIDKey contextKey = "request_id"

// RequestID 中间件 - 生成唯一request_id并通过context传递
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 生成唯一 request_id
		requestID := generateRequestID()

		// 将request_id存入context中
		ctx := context.WithValue(r.Context(),RequestIDKey,requestID)
		r = r.WithContext(ctx)

		//在响应头中也返回request_id
		w.Header().Set("X-Request-ID",requestID) 

		//记录日志
		log.Printf("[RequestID: %s] New request: %s %s", requestID, r.Method, r.URL.Path)

		//调用下一个处理器
		next.ServeHTTP(w,r)
	})
}

// generateRequestID 生成唯一 request_id
func generateRequestID() string {
	//方法1 ：时间戳 + 随机数
	timestamp := time.Now().UnixMilli() //毫秒
	randomBytes := make([]byte,4)
	
	// 读取随机数，如果失败则使用时间戳作为随机源
	n, err := rand.Read(randomBytes)
	if err != nil || n != 4 {
		// 降级方案：使用纳秒时间戳的后8位作为随机数
		nanoTime := time.Now().UnixNano()
		randomBytes = []byte{
			byte(nanoTime >> 24),
			byte(nanoTime >> 16),
			byte(nanoTime >> 8),
			byte(nanoTime),
		}
	}
	
	randomHex := hex.EncodeToString(randomBytes)//hex.EncodeToString() - 将字节转换为十六进制字符串
												//例如：[0x3a, 0xf2, 0x1c, 0x8d] → "3af21c8d"

	return fmt.Sprintf("%d-%s", timestamp, randomHex)
}

// GetRequestID 从 context 中获取 request_id
func GetRequestID(ctx context.Context) string{
	if requestID,ok := ctx.Value(RequestIDKey).(string);ok{
		return requestID
	}
	return ""
}