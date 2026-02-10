package core

import (
	"context"

	"github.com/voltavpn/volta-client/internal/api"
	"github.com/voltavpn/volta-client/internal/authlink"
)

func ValidateAccessInput(raw string) (statusMessage string, ok bool) {
	normalized := authlink.NormalizeInput(raw)
	if normalized == "" {
		return "Пожалуйста, введите ключ доступа или ссылку.", false
	}

	_, tokenOK := authlink.ExtractToken(normalized)
	if !tokenOK {
		return "Неверный ключ или ссылка.", false
	}

	return "Формат ключа принят. Продолжаем…", true
}

type ActivateResult struct {
	SessionToken string
	VPNProfile   string
	ProfileURL   string
}

func ActivateAccess(ctx context.Context, client api.APIClient, raw string) (ActivateResult, string, bool) {
	var empty ActivateResult

	normalized := authlink.NormalizeInput(raw)
	if normalized == "" {
		return empty, "Пожалуйста, введите ключ доступа или ссылку.", false
	}

	token, tokenOK := authlink.ExtractToken(normalized)
	if !tokenOK {
		return empty, "Неверный ключ или ссылка.", false
	}

	if client == nil {
		return empty, "Сервис временно недоступен. Повторите попытку позже.", false
	}

	resp, err := client.Activate(ctx, token)
	if err != nil || resp == nil {
		// Сообщение умышленно общее, без раскрытия деталей сетевой ошибки.
		return empty, "Не удалось связаться с сервером. Повторите попытку позже.", false
	}

	result := ActivateResult{
		SessionToken: resp.SessionToken,
		VPNProfile:   resp.VPNProfile,
		ProfileURL:   resp.ProfileURL,
	}

	return result, "Ключ подтверждён. Готовим подключение…", true
}
