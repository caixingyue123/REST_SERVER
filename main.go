package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"restServer/middleware"
	"restServer/model"
	"restServer/response"
)

// ============ Day 1-2: 基础路由 ============

// 任务 1.1: 最简单的 HTTP 服务
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

// 任务 1.2: 多路由处理
func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Go HTTP Server!")
}

// 任务 1.3: 处理不同 HTTP 方法
func methodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintf(w, "GET request received")
	case http.MethodPost:
		fmt.Fprintf(w, "POST request received")
	case http.MethodPut:
		fmt.Fprintf(w, "PUT request received")
	case http.MethodDelete:
		fmt.Fprintf(w, "DELETE request received")
	default:
		// 返回 405 Method Not Allowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
	}
}

// ============ Day 3: JSON 编解码 + 统一返回格式 ============

// 任务 2.1: JSON 响应 - 使用统一返回格式
func jsonResponseHandler(w http.ResponseWriter, r *http.Request) {
	// data := map[string]any{
	// 	"name": "张三",
	// 	"age":  25,
	// 	"city": "北京",
	// }
	user := &model.User{
		Name:     "李四",
		Password: "afds8979",
	}
	response.Success(w, &user)
}

// 任务 2.2 & 2.3: JSON 请求解析 + 参数校验
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// 只允许 POST 方法
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, 405, "Method not allowed")
		return
	}

	// 解析 JSON 请求体
	var user model.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		response.BadRequest(w, "JSON 格式错误: "+err.Error())
		return
	}
	defer r.Body.Close()

	//参数校验
	if err := user.Validate(); err != nil {
		response.Error(w, http.StatusForbidden, 400, "参数校验错误")
	}

	//返回成功响应
	response.Success(w, map[string]any{
		"message": "用户创建成功",
		"user": map[string]any{
			"username": user.Name,
			"email":    user.Email,
		},
	})

}

// 测试错误响应
func errorTestHandler(w http.ResponseWriter, r *http.Request) {
	response.InternalError(w, "这是一个模拟的服务器错误")
}

// 测试 panic 恢复
func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("这是一个故意触发的panic")
}

// 需要鉴权的接口
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	response.Success(w, map[string]any{
		"message":    "这是一个受保护的资源",
		"request_id": requestID,
	})
}

// 测试request_id
func requestIDTestHandler(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetRequestID(r.Context())
	response.Success(w, map[string]any{
		"message":    "Request ID 测试",
		"request_id": requestID,
	})
}

// ============ Day 5: Context 超时控制 ============

// 任务 4.1: 模拟慢接口，测试超时效果
func slowHandler(w http.ResponseWriter, r *http.Request) {
	//设置5秒的超时
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	//模拟耗时操作
	select {
	case <-time.After(time.Second * 3):
		//正常完成
		response.Success(w, map[string]any{
			"message":  "慢接口处理完成",
			"duration": "3秒",
		})
	case <-ctx.Done():
		//超时或取消
		response.Error(w, http.StatusRequestTimeout, 40801, "请求超时")
		return
	}
}

// 在业务逻辑里面使用context
func contextBusinessHandler(w http.ResponseWriter, r *http.Request) {
	//设置2秒的超时
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()

	//模拟数据库查询
	result, err := simulateDBQuery(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			response.Error(w, http.StatusRequestTimeout, 40801, "数据库查询超时")
		} else {
			response.InternalError(w, "数据库查询失败"+err.Error())
		}
		return
	}
	response.Success(w, result)
}

func simulateDBQuery(ctx context.Context) (map[string]any, error) {
	//创建一个通道来接受结果
	resultChan := make(chan map[string]any, 1)
	errorChan := make(chan error, 1)

	//启动goroutine模拟数据库查询
	go func() {
		//模拟查询耗时1秒
		time.Sleep(time.Millisecond * 1500)

		//返回模拟数据
		resultChan <- map[string]any{
			"user": []map[string]any{
				{"id": 1, "name": "张三", "email": "zhangsan@example.com"},
				{"id": 2, "name": "李四", "email": "lisi@example.com"},
			},
			"total": 2,
		}
	}()
	//等待结果或者超时
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// 超时测试接口 - 故意超时
func timeoutTestHandler(w http.ResponseWriter, r *http.Request) {
	//设置1秒超时，但是模拟3秒操作
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	select {
	case <-time.After(time.Second * 3):
		response.Success(w, map[string]any{
			"message": "不应该看到这个消息",
		})
	case <-ctx.Done():
		response.Error(w, http.StatusRequestTimeout, 40801, "操作超时")
		return
	}
}

func main() {
	// 创建一个新的 ServeMux
	mux := http.NewServeMux()

	// Day 1-2 路由
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/method", methodHandler)

	// Day 3 路由
	mux.HandleFunc("/api/json", jsonResponseHandler)
	mux.HandleFunc("/api/user", createUserHandler)
	mux.HandleFunc("/api/error", errorTestHandler)

	// Day 4 路由 - 公开接口
	mux.HandleFunc("/api/panic", panicHandler)
	mux.HandleFunc("/api/requestid", requestIDTestHandler)

	// Day 5 路由 - Context超时控制
	mux.HandleFunc("/api/slow", slowHandler)
	mux.HandleFunc("/api/context-business", contextBusinessHandler)
	mux.HandleFunc("/api/timeout-test", timeoutTestHandler)

	// 应用基础中间件链到整个 mux
	// 链式顺序: Recovery -> RequestID -> Logger -> mux
	handler := middleware.Recovery(
		middleware.RequestID(
			middleware.Logger(mux),
		),
	)

	// 创建需要鉴权的路由处理器
	protectedHandlerFunc := http.HandlerFunc(protectedHandler)
	protectedWithAuth := middleware.Auth(protectedHandlerFunc)

	// 注册受保护的路由（已经包含了基础中间件，这里只加鉴权）
	authMux := http.NewServeMux()
	authMux.Handle("/api/protected", protectedWithAuth)

	// 合并所有路由
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 如果是受保护的路由，先应用基础中间件，再应用鉴权
		if r.URL.Path == "/api/protected" {
			middleware.Recovery(
				middleware.RequestID(
					middleware.Logger(
						middleware.Auth(protectedHandlerFunc),
					),
				),
			).ServeHTTP(w, r)
		} else {
			// 否则使用普通处理器
			handler.ServeHTTP(w, r)
		}
	})
	server := &http.Server{
		Addr:         ":8080",
		Handler:      finalHandler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 60,
	}

	// 启动服务器
	go func() {
		fmt.Println("Server is running on http://localhost:8080")
		fmt.Println("\n=== Day 4 中间件测试接口 ===")
		fmt.Println("公开接口:")
		fmt.Println("  GET  /api/requestid - 测试 Request ID 中间件")
		fmt.Println("  GET  /api/panic     - 测试 Recovery 中间件")
		fmt.Println("\n受保护接口 (需要 Token):")
		fmt.Println("  GET  /api/protected - 测试鉴权中间件")
		fmt.Println("\n测试命令:")
		fmt.Println("  # 测试 Request ID")
		fmt.Println("  curl http://localhost:8080/api/requestid")
		fmt.Println("\n  # 测试 Panic Recovery")
		fmt.Println("  curl http://localhost:8080/api/panic")
		fmt.Println("\n  # 测试鉴权失败")
		fmt.Println("  curl http://localhost:8080/api/protected")
		fmt.Println("\n  # 测试鉴权成功")
		fmt.Println(`  curl -H "Authorization: Bearer my-secret-token-123456" http://localhost:8080/api/protected`)

		// 启动服务器 - 使用合并后的 handler
		if err := http.ListenAndServe(":8080", finalHandler); err != nil {
			fmt.Printf("Server failed to start: %v\n", err)
		}
	}()

	// 任务 5.2: 实现优雅退出
	// 创建一个通道来接收系统信号
	quit := make(chan os.Signal, 1)

	// 监听 SIGINT (Ctrl+C) 和 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待信号
	<-quit
	log.Println("收到退出信号，开始优雅关闭服务器...")

	// 设置关闭超时时间（30秒）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("服务器强制关闭: %v", err)
	} else {
		log.Println("服务器已优雅关闭！")
	}
}
