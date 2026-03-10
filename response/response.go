package response

import (
	"encoding/json"
	"net/http"
)

// Response 统一返回结构体
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// WriteJSON 封装 JSON 响应函数
func WriteJSON(w http.ResponseWriter, statusCode int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

// Success 成功响应
func Success(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// ERROR 失败响应 : statusCode是状态码，不同的错误状态码不同，code是业务错误码，更加细节的代码
func Error(w http.ResponseWriter, statusCode int, code int, message string) {
	WriteJSON(w, statusCode, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400 错误
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, 40001, message)
}

// InternalError 500 错误
func InternalError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, 50001, message)
}
