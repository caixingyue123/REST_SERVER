package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restServer/model"
	"restServer/response"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 内存存储 todos
var (
	todos     = make(map[int]*model.Todo)
	todoIDSeq = 1
	todoMutex sync.RWMutex
)

// CreateTodoRequest 创建 Todo 请求
type CreateTodoRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdateTodoRequest 更新 Todo 请求
type UpdateTodoRequest struct {
	Title     *string `json:"title,omitempty"`
	Content   *string `json:"content,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}

// ListTodos 获取 Todo 列表
func ListTodos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, 405, "Method not allowed")
		return
	}

	todoMutex.RLock()
	defer todoMutex.RUnlock()

	list := make([]*model.Todo, 0, len(todos))
	for _, todo := range todos {
		list = append(list, todo)
	}

	response.Success(w, map[string]any{
		"todos": list,
		"total": len(list),
	})
}

// CreateTodo 创建 Todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, 405, "Method not allowed")
		return
	}

	// 检查 Content-Type
	if ct := r.Header.Get("Content-Type"); ct != "" && !strings.Contains(ct, "application/json") {
		response.BadRequest(w, "Content-Type 必须为 application/json")
		return
	}

	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON 格式错误: "+err.Error())
		return
	}
	defer r.Body.Close()

	// 参数校验
	if req.Title == "" {
		response.BadRequest(w, "标题不能为空")
		return
	}
	if len(req.Title) > 100 {
		response.BadRequest(w, "标题长度不能超过100个字符")
		return
	}
	if len(req.Content) > 1000 {
		response.BadRequest(w, "内容长度不能超过1000个字符")
		return
	}

	todoMutex.Lock()
	defer todoMutex.Unlock()

	now := time.Now()
	todo := &model.Todo{
		ID:        todoIDSeq,
		Title:     req.Title,
		Content:   req.Content,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	todos[todoIDSeq] = todo
	todoIDSeq++

	response.Success(w, map[string]any{
		"message": "创建成功",
		"todo":    todo,
	})
}

// UpdateTodo 更新 Todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.Error(w, http.StatusMethodNotAllowed, 405, "Method not allowed")
		return
	}

	// 从 URL 中提取 ID: /api/todos/:id
	id, err := extractIDFromPath(r.URL.Path, "/api/todos/")
	if err != nil {
		response.BadRequest(w, "无效的 ID")
		return
	}

	// 检查 Content-Type
	if ct := r.Header.Get("Content-Type"); ct != "" && !strings.Contains(ct, "application/json") {
		response.BadRequest(w, "Content-Type 必须为 application/json")
		return
	}

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "JSON 格式错误: "+err.Error())
		return
	}
	defer r.Body.Close()

	// 检查是否提供了至少一个更新字段
	if req.Title == nil && req.Content == nil && req.Completed == nil {
		response.BadRequest(w, "至少需要提供一个更新字段")
		return
	}

	todoMutex.Lock()
	defer todoMutex.Unlock()

	todo, exists := todos[id]
	if !exists {
		response.Error(w, http.StatusNotFound, 40401, "Todo 不存在")
		return
	}

	// 更新字段
	if req.Title != nil {
		if *req.Title == "" {
			response.BadRequest(w, "标题不能为空")
			return
		}
		todo.Title = *req.Title
	}
	if req.Content != nil {
		todo.Content = *req.Content
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	todo.UpdatedAt = time.Now()

	response.Success(w, map[string]any{
		"message": "更新成功",
		"todo":    todo,
	})
}

// DeleteTodo 删除 Todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.Error(w, http.StatusMethodNotAllowed, 405, "Method not allowed")
		return
	}

	// 从 URL 中提取 ID
	id, err := extractIDFromPath(r.URL.Path, "/api/todos/")
	if err != nil {
		response.BadRequest(w, "无效的 ID")
		return
	}

	todoMutex.Lock()
	defer todoMutex.Unlock()

	if _, exists := todos[id]; !exists {
		response.Error(w, http.StatusNotFound, 40401, "Todo 不存在")
		return
	}

	delete(todos, id)

	response.Success(w, map[string]any{
		"message": "删除成功",
		"id":      id,
	})
}

// extractIDFromPath 从路径中提取 ID
func extractIDFromPath(path, prefix string) (int, error) {
	idStr := strings.TrimPrefix(path, prefix)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	if id <= 0 {
		return 0, fmt.Errorf("ID 必须为正整数")
	}
	return id, nil
}
