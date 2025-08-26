package apperror

import "net/http"

type AppError struct {
	Code    int    // HTTP код (например, 400, 404, 500)
	Message string // Человекочитаемое сообщение
	Err     error  // Оригинальная ошибка (wrap)
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, msg string, err error) *AppError {
	return &AppError{Code: code, Message: msg, Err: err}
}

func BadRequest(msg string, err error) *AppError {
	return New(http.StatusBadRequest, msg, err)
}

func NotFound(msg string, err error) *AppError {
	return New(http.StatusNotFound, msg, err)
}

func Internal(msg string, err error) *AppError {
	return New(http.StatusInternalServerError, msg, err)
}
