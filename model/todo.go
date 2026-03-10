package model

import (
	"fmt"
	"time"
)

// Todo 待办事项模型
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 参数校验
func (t *Todo) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("标题不能为空")
	}
	if len(t.Title) > 100 {
		return fmt.Errorf("标题长度不能超过100个字符")
	}
	return nil
}
