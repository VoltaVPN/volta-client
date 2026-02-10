package core

import (
	"github.com/voltavpn/volta-client/internal/authlink"
)

// ValidateAccessInput делает безопасную локальную проверку введённого значения.
// Важно: функция не логирует и не печатает введённые данные.
func ValidateAccessInput(raw string) (statusMessage string, ok bool) {
	normalized := authlink.NormalizeInput(raw)
	if normalized == "" {
		return "Пожалуйста, введите ключ доступа или ссылку.", false
	}

	_, tokenOK := authlink.ExtractToken(normalized)
	if !tokenOK {
		// Сообщение нарочно общее, без подсказок для перебора.
		return "Неверный ключ или ссылка.", false
	}

	// На этом этапе формат токена и домен считаем валидными.
	// Следующие шаги (аутентификация, подключение к VPN) реализуются отдельно.
	return "Формат ключа принят. Продолжаем…", true
}
