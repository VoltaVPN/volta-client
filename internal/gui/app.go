package gui

import (
	"context"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/voltavpn/volta-client/internal/api"
	"github.com/voltavpn/volta-client/internal/core"
)

func Run() {
	application := app.New()
	application.Settings().SetTheme(NewVoltaTheme())
	window := application.NewWindow("VoltaVPN")

	apiClient, err := api.NewClientFromEnv()
	if err != nil {
		showErrorScreen(window, "Сервис временно недоступен. Повторите попытку позже.")
		window.ShowAndRun()
		return
	}

	// Dev shortcut: для теста UI главного экрана без логина/активации.
	// Важно: используем безопасный stub и не показываем/не логируем секреты.
	if strings.TrimSpace(os.Getenv("VOLTA_DEV_SKIP_LOGIN")) == "1" {
		showMainScreen(window, core.ActivateResult{
			// SessionToken/VPNProfile/ProfileURL намеренно пустые.
		})
	} else {
		showLoginScreen(window, apiClient)
	}

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

		showMainScreen(window, result)
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

func showActivatingScreen(window fyne.Window, result core.ActivateResult) {
	titleLabel := widget.NewLabelWithStyle(
		"VoltaVPN",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	statusLabel := widget.NewLabelWithStyle(
		"Activating…",
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

func showMainScreen(window fyne.Window, result core.ActivateResult) {
	titleLabel := widget.NewLabelWithStyle(
		"VoltaVPN",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	statusBadge := newStatusBadge("Disconnected", StatusDisconnectedColor())

	statusRow := container.NewHBox(
		widget.NewLabel("Status"),
		layout.NewSpacer(),
		statusBadge,
	)

	connectButton := widget.NewButton("Connect", func() {
		statusBadge.SetText("Coming soon")
	})
	connectButton.Importance = widget.HighImportance

	settingsButton := widget.NewButton("Settings", func() {
		showSettingsScreen(window, result)
	})

	profileName, profileRegion := safeProfileView(result)
	profileContent := container.NewVBox(
		makeLabelRow("Имя", profileName),
		makeLabelRow("Регион", profileRegion),
	)
	activeProfileCard := widget.NewCard("Активный профиль", "", profileContent)

	uploadLabel := canvas.NewText("↑ Upload 0 B/s", AccentColor())
	uploadLabel.TextSize = 12
	downloadLabel := canvas.NewText("↓ Download 0 B/s", AccentColor())
	downloadLabel.TextSize = 12
	statsRow := container.NewHBox(uploadLabel, layout.NewSpacer(), downloadLabel)

	content := container.NewVBox(
		titleLabel,
		layout.NewSpacer(),
		statusRow,
		layout.NewSpacer(),
		connectButton,
		settingsButton,
		layout.NewSpacer(),
		activeProfileCard,
		statsRow,
	)

	padded := container.NewPadded(content)
	window.SetContent(container.NewCenter(padded))
}

func makeLabelRow(label, value string) *fyne.Container {
	return container.NewHBox(
		widget.NewLabel(label),
		layout.NewSpacer(),
		widget.NewLabel(value),
	)
}

func showSettingsScreen(window fyne.Window, result core.ActivateResult) {
	titleLabel := widget.NewLabelWithStyle(
		"Настройки",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	backLink := widget.NewHyperlink("← Назад", nil)
	backLink.OnTapped = func() {
		showMainScreen(window, result)
	}

	placeholder := widget.NewLabel("Coming soon")

	content := container.NewVBox(
		titleLabel,
		layout.NewSpacer(),
		placeholder,
		layout.NewSpacer(),
		backLink,
	)

	padded := container.New(
		layout.NewVBoxLayout(),
		container.NewPadded(content),
	)

	window.SetContent(container.NewCenter(padded))
}

func safeProfileView(result core.ActivateResult) (name string, region string) {
	// Строгие требования безопасности:
	// - Никогда не используем result.SessionToken / result.VPNProfile / result.ProfileURL в UI.
	// - Не пытаемся парсить конфиги или URL, чтобы случайно не показать внутренние детали.
	//
	// Поэтому сейчас отображаем только безопасные заглушки.
	return "Active profile", "Unknown"
}
