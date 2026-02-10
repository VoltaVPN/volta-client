package gui

import (
	"context"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/voltavpn/volta-client/internal/api"
	"github.com/voltavpn/volta-client/internal/core"
)

func Run() {
	application := app.New()
	window := application.NewWindow("VoltaVPN")

	apiClient, err := api.NewClientFromEnv()
	if err != nil {
		showErrorScreen(window, "Сервис временно недоступен. Повторите попытку позже.")
		window.ShowAndRun()
		return
	}

	showLoginScreen(window, apiClient)

	window.Resize(fyne.NewSize(420, 260))
	window.CenterOnScreen()

	window.ShowAndRun()
}

func showLoginScreen(window fyne.Window, apiClient api.APIClient) {
	titleLabel := widget.NewLabelWithStyle(
		"VoltaVPN",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	subtitleLabel := widget.NewLabelWithStyle(
		"Войдите с помощью ключа доступа или ссылки",
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)

	accessLabel := widget.NewLabel("Ключ доступа или ссылка")

	accessInputEntry := widget.NewEntry()
	accessInputEntry.SetPlaceHolder("https://vvpn.io/...")

	statusLabel := widget.NewLabel("")
	statusLabel.Wrapping = fyne.TextWrapWord

	var lastAttempt time.Time

	continueButton := widget.NewButton("Продолжить", nil)
	continueButton.OnTapped = func() {
		now := time.Now()
		if !lastAttempt.IsZero() && now.Sub(lastAttempt) < 500*time.Millisecond {
			return
		}
		lastAttempt = now

		continueButton.Disable()
		accessInputEntry.Disable()
		statusLabel.SetText("Проверяем ключ…")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		result, statusMessage, ok := core.ActivateAccess(ctx, apiClient, accessInputEntry.Text)
		statusLabel.SetText(statusMessage)

		if !ok {
			accessInputEntry.Enable()
			if strings.TrimSpace(accessInputEntry.Text) != "" {
				continueButton.Enable()
			}
			return
		}

		showConnectingScreen(window, result)
	}
	continueButton.Importance = widget.HighImportance
	continueButton.Disable()

	accessInputEntry.OnChanged = func(text string) {
		if strings.TrimSpace(text) == "" {
			continueButton.Disable()
			return
		}
		continueButton.Enable()
	}

	form := container.NewVBox(
		titleLabel,
		subtitleLabel,
		layout.NewSpacer(),
		accessLabel,
		accessInputEntry,
		continueButton,
		statusLabel,
	)

	formContainer := container.New(
		layout.NewVBoxLayout(),
		container.NewPadded(form),
	)

	centeredContent := container.NewCenter(formContainer)

	window.SetContent(centeredContent)
}

func showConnectingScreen(window fyne.Window, result core.ActivateResult) {
	titleLabel := widget.NewLabelWithStyle(
		"VoltaVPN",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	statusLabel := widget.NewLabelWithStyle(
		"Подключение…",
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)
	statusLabel.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		titleLabel,
		layout.NewSpacer(),
		statusLabel,
		layout.NewSpacer(),
	)

	padded := container.New(
		layout.NewVBoxLayout(),
		container.NewPadded(content),
	)

	center := container.NewCenter(padded)
	window.SetContent(center)
}

func showErrorScreen(window fyne.Window, message string) {
	titleLabel := widget.NewLabelWithStyle(
		"VoltaVPN",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	statusLabel := widget.NewLabelWithStyle(
		message,
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)
	statusLabel.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		titleLabel,
		layout.NewSpacer(),
		statusLabel,
		layout.NewSpacer(),
	)

	padded := container.New(
		layout.NewVBoxLayout(),
		container.NewPadded(content),
	)

	center := container.NewCenter(padded)
	window.SetContent(center)
}
