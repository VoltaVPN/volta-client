package authlink

import (
	"net/url"
	"strings"
)

const (
	allowedHostSuffix = ".vvpn.io"
	allowedHostExact  = "vvpn.io"

	minTokenLen = 32
	maxTokenLen = 512
)

// NormalizeInput приводит пользовательский ввод к каноничному виду
// (обрезает пробелы по краям). Важно не модифицировать само значение токена.
func NormalizeInput(s string) string {
	return strings.TrimSpace(s)
}

// ExtractToken извлекает токен как из полной ссылки, так и из "голого" токена.
// Возвращает ok=false, если значение не подходит по формату или домену.
func ExtractToken(s string) (token string, ok bool) {
	normalized := NormalizeInput(s)
	if normalized == "" {
		return "", false
	}

	// Вариант 1: похоже на URL с протоколом.
	if strings.HasPrefix(normalized, "http://") || strings.HasPrefix(normalized, "https://") {
		u, err := url.Parse(normalized)
		if err != nil {
			return "", false
		}

		host := strings.ToLower(u.Hostname())
		if !isAllowedHost(host) {
			return "", false
		}

		// Дополнительные параметры/фрагменты не разрешаем.
		if u.RawQuery != "" || u.Fragment != "" {
			return "", false
		}

		path := strings.Trim(u.EscapedPath(), "/")
		if path == "" || strings.Contains(path, "/") {
			return "", false
		}

		token = path
	} else {
		// Вариант 2: пользователь вставил только токен.
		token = normalized
	}

	if !ValidateTokenFormat(token) {
		return "", false
	}

	return token, true
}

// ValidateTokenFormat проверяет формат токена без привязки к домену/URL.
// Требования:
//   - длина в разумном диапазоне (minTokenLen..maxTokenLen)
//   - только URL-safe base64-like символы [A-Za-z0-9_-]
func ValidateTokenFormat(token string) bool {
	if len(token) < minTokenLen || len(token) > maxTokenLen {
		return false
	}

	for i := 0; i < len(token); i++ {
		c := token[i]
		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '-' || c == '_' {
			continue
		}
		return false
	}

	return true
}

func isAllowedHost(host string) bool {
	if host == allowedHostExact {
		return true
	}
	if strings.HasSuffix(host, allowedHostSuffix) {
		return true
	}
	return false
}
