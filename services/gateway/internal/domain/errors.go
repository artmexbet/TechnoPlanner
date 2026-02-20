package domain

import "errors"

// Domain errors - используются для изоляции слоя БД от слоя сервисов
var (
	// ErrNotFound - ресурс не найден
	ErrNotFound = errors.New("not found")
	// ErrForbidden - доступ запрещен
	ErrForbidden = errors.New("forbidden")
	// ErrInvalidInput - невалидные входные данные
	ErrInvalidInput = errors.New("invalid input")
	// ErrConflict - конфликт данных
	ErrConflict = errors.New("conflict")
)
