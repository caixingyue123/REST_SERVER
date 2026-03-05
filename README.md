# Go 语言学习 - 第 2 周：HTTP 服务从 0 到 1

## 学习目标

能够自己写一个最小但"像样"的 HTTP 服务，掌握：
- 路由、handler、中间件（日志、recover、鉴权占位）
- JSON 编解码、参数校验（先手写也行）
- 统一返回格式：code/message/data
- 超时控制：context.WithTimeout
- 优雅退出：http.Server + Shutdown

## 本周产出

- 一个 REST 服务
- 用户注册/登录（先假数据也行）
- 2~3 个 CRUD 接口
- 带 request_id 的日志（哪怕先用简单中间件）

---

## Day 1-2: net/http 基础 + 路由与 Handler

### 任务 1.1: 最简单的 HTTP 服务
- [ ] 使用 `net/http` 创建一个 Hello World 服务
- [ ] 理解 `http.HandleFunc` 和 `http.ListenAndServe`
- [ ] 测试：浏览器访问 `localhost:8080`

### 任务 1.2: 多路由处理
- [ ] 实现 3-4 个不同路径的 handler（如 `/`, `/ping`, `/health`）
- [ ] 理解 `http.ServeMux` 的路由匹配规则
- [ ] 学习 `http.Handler` 接口

### 任务 1.3: 处理不同 HTTP 方法
- [ ] 在 handler 中区分 GET/POST/PUT/DELETE
- [ ] 返回 405 Method Not Allowed

---

## Day 3: JSON 编解码 + 统一返回格式

### 任务 2.1: JSON 响应
- [ ] 定义统一返回结构体：`Response{Code int, Message string, Data interface{}}`
- [ ] 封装 `WriteJSON` 函数
- [ ] 实现成功/失败响应的辅助函数

### 任务 2.2: JSON 请求解析
- [ ] 使用 `json.Decoder` 解析请求体
- [ ] 处理解析错误（400 Bad Request）
- [ ] 实现一个简单的 POST 接口测试

### 任务 2.3: 参数校验
- [ ] 手写基础校验逻辑（非空、长度、格式）
- [ ] 封装 `Validate` 函数
- [ ] 返回友好的错误提示

---

## Day 4: 中间件机制

### 任务 3.1: 理解中间件模式
- [ ] 学习中间件的洋葱模型
- [ ] 实现第一个中间件：日志中间件（记录请求方法、路径、耗时）

### 任务 3.2: Request ID 中间件
- [ ] 生成唯一 request_id（可用 `uuid` 或简单的时间戳+随机数）
- [ ] 通过 `context.Context` 传递 request_id
- [ ] 在日志中打印 request_id

### 任务 3.3: Recover 中间件
- [ ] 捕获 panic，避免服务崩溃
- [ ] 返回 500 错误
- [ ] 记录错误堆栈

### 任务 3.4: 鉴权中间件占位
- [ ] 实现简单的 Token 校验逻辑（先硬编码 token）
- [ ] 从 Header 中读取 `Authorization`
- [ ] 返回 401 Unauthorized

---

## Day 5: Context 超时控制

### 任务 4.1: 理解 Context
- [ ] 学习 `context.Context` 的作用
- [ ] 使用 `context.WithTimeout` 设置请求超时
- [ ] 模拟慢接口，测试超时效果

### 任务 4.2: 在业务逻辑中使用 Context
- [ ] 在 handler 中传递 context
- [ ] 监听 `ctx.Done()` 提前退出

---

## Day 6: 优雅退出

### 任务 5.1: 使用 http.Server
- [ ] 从 `http.ListenAndServe` 改为 `http.Server` 结构体
- [ ] 配置 `ReadTimeout`、`WriteTimeout`

### 任务 5.2: 实现优雅退出
- [ ] 监听系统信号（`os.Signal`，如 SIGINT、SIGTERM）
- [ ] 调用 `server.Shutdown(ctx)` 等待请求处理完成
- [ ] 设置 shutdown 超时时间

---

## Day 7: 综合实战 - REST 服务

### 任务 6.1: 用户注册接口
- [ ] `POST /api/register`
- [ ] 接收 `username` 和 `password`（先用内存 map 存储）
- [ ] 参数校验 + 返回统一格式

### 任务 6.2: 用户登录接口
- [ ] `POST /api/login`
- [ ] 校验用户名密码
- [ ] 返回简单 token（可以是随机字符串）

### 任务 6.3: CRUD 接口（以 Todo 为例）
- [ ] `GET /api/todos` - 列表查询
- [ ] `POST /api/todos` - 创建
- [ ] `PUT /api/todos/:id` - 更新
- [ ] `DELETE /api/todos/:id` - 删除
- [ ] 数据先用内存 slice 存储

### 任务 6.4: 应用所有中间件
- [ ] 串联：日志 → request_id → recover → 鉴权（部分接口）
- [ ] 测试完整请求链路

---

## 加分项（时间充裕可做）

### 任务 7.1: 引入 Gin 框架
- [ ] 安装 Gin：`go get -u github.com/gin-gonic/gin`
- [ ] 用 Gin 重写上面的接口
- [ ] 对比 net/http 和 Gin 的区别

### 任务 7.2: 使用 Gin 中间件
- [ ] 使用 Gin 内置的 `Logger()`、`Recovery()`
- [ ] 自定义 request_id 中间件
- [ ] 学习 Gin 的参数绑定和校验（`ShouldBindJSON`、`binding` tag）

---

## 学习建议

1. **每天写代码**：不要只看不练，每个任务都要自己敲一遍
2. **测试驱动**：用 Postman 或 curl 测试每个接口
3. **记录笔记**：遇到的坑和解决方案记下来
4. **代码提交**：每完成一个任务就 git commit，方便回顾
5. **循序渐进**：先用 net/http 打好基础，再用 Gin 体会框架的便利

---

## 测试命令示例

```bash
# 测试 GET 请求
curl http://localhost:8080/ping

# 测试 POST 请求
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'

# 带 Token 的请求
curl http://localhost:8080/api/todos \
  -H "Authorization: Bearer your-token-here"
```

---

## 项目结构建议

```
.
├── main.go              # 入口文件
├── handler/             # 处理器
│   ├── user.go
│   └── todo.go
├── middleware/          # 中间件
│   ├── logger.go
│   ├── recovery.go
│   ├── requestid.go
│   └── auth.go
├── model/              # 数据模型
│   ├── user.go
│   └── todo.go
└── response/           # 统一响应
    └── response.go
```
