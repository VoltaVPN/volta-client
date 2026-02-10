package core

import "strings"

// ValidateAccessInput делает простую локальную проверку введённого значения.
func ValidateAccessInput(raw string) (statusMessage string, ok bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "Пожалуйста, введите ключ доступа или ссылку.", false
	}

	return "Формат введённых данных принят (локальная проверка).", true
}
