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
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/voltavpn/volta-client/internal/api"
	"github.com/voltavpn/volta-client/internal/core"
	"github.com/voltavpn/volta-client/internal/settings"
	"github.com/voltavpn/volta-client/internal/ui/components"
)

func Run() {
	application := app.New()
	application.Settings().SetTheme(NewVoltaTheme())
	window := application.NewWindow("VoltaVPN")

	appSettings := settings.LoadOrDefault()

	apiClient, err := api.NewClientFromEnv()
	if err != nil {
		showErrorScreen(window, "Сервис временно недоступен. Повторите попытку позже.")
		window.ShowAndRun()
		return
	}

	// VOLTA_DEV_SKIP_LOGIN допускается только в dev-окружении.
	if isDevEnvironment() && strings.TrimSpace(os.Getenv("VOLTA_DEV_SKIP_LOGIN")) == "1" {
		showMainScreen(window, core.ActivateResult{}, &appSettings)
	} else {
		showLoginScreen(window, apiClient, &appSettings)
	}

	window.Resize(fyne.NewSize(560, 560))
	window.CenterOnScreen()
	window.ShowAndRun()
}

func showLoginScreen(window fyne.Window, apiClient api.APIClient, appSettings *settings.Settings) {
	titleLabel := canvas.NewText("VoltaVPN", components.ColorText())
	titleLabel.TextSize = components.TextHeadline
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter

	promptLabel := canvas.NewText("Введите ссылку доступа", components.ColorTextMuted())
	promptLabel.TextSize = components.BodyTextSize
	promptLabel.Alignment = fyne.TextAlignCenter

	accessInputEntry := widget.NewEntry()
	accessInputEntry.SetPlaceHolder("https://...")
	accessInputEntry.Wrapping = fyne.TextWrapOff

	continueButton := components.NewPrimaryButton("Продолжить", nil)

	var lastAttempt time.Time
	continueButton.OnTapped = func() {
		accessURL := strings.TrimSpace(accessInputEntry.Text)
		if accessURL == "" {
			dialog.ShowInformation("Ввод ссылки", "Введите ссылку доступа.", window)
			window.Canvas().Focus(accessInputEntry)
			return
		}

		now := time.Now()
		if !lastAttempt.IsZero() && now.Sub(lastAttempt) < 500*time.Millisecond {
			return
		}
		lastAttempt = now

		continueButton.Disable()
		continueButton.SetText("Проверяем...")
		accessInputEntry.Disable()

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		result, _, ok := core.ActivateAccess(ctx, apiClient, accessURL)
		if !ok {
			accessInputEntry.Enable()
			continueButton.SetText("Продолжить")
			continueButton.Enable()
			dialog.ShowInformation("Ошибка", "Не удалось активировать доступ. Проверьте ссылку и попробуйте снова.", window)
			return
		}

		showMainScreen(window, result, appSettings)
	}

	privacyCaption := canvas.NewText("Ключ не сохраняется в открытом виде", components.ColorTextMuted())
	privacyCaption.TextSize = components.CaptionTextSize
	privacyCaption.Alignment = fyne.TextAlignCenter

	form := container.NewVBox(
		titleLabel,
		components.NewVSpacer(components.Spacing8),
		promptLabel,
		components.NewVSpacer(components.Spacing8),
		accessInputEntry,
		components.NewVSpacer(components.Spacing8),
		continueButton,
		components.NewVSpacer(components.Spacing8),
		privacyCaption,
	)

	content := container.NewBorder(
		nil,
		nil,
		components.NewHSpacer(components.Spacing16),
		components.NewHSpacer(components.Spacing16),
		container.NewVBox(
			layout.NewSpacer(),
			form,
			layout.NewSpacer(),
		),
	)
	window.SetContent(content)
}

func showErrorScreen(window fyne.Window, message string) {
	titleLabel := canvas.NewText("VoltaVPN", components.ColorText())
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.TextSize = components.TitleTextSize
	titleLabel.Alignment = fyne.TextAlignCenter

	statusLabel := canvas.NewText(message, components.ColorTextMuted())
	statusLabel.Alignment = fyne.TextAlignCenter
	statusLabel.TextSize = components.BodyTextSize

	card := components.NewCardWithPadding(
		container.NewVBox(
			titleLabel,
			components.NewVSpacer(components.Spacing16),
			statusLabel,
		),
		components.Spacing24,
		components.Spacing24,
	)

	window.SetContent(container.NewCenter(container.NewPadded(card)))
}

func showMainScreen(window fyne.Window, result core.ActivateResult, appSettings *settings.Settings) {
	titleLabel := canvas.NewText("VoltaVPN", components.ColorText())
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.TextSize = components.TextTitle

	statusLabel := canvas.NewText("Status: Disconnected", components.ColorStatusDisconnected())
	statusLabel.TextSize = components.TextBody

	userIDLabel := canvas.NewText("ID: -", components.ColorTextMuted())
	userIDLabel.TextSize = components.TextBody
	if hasAnyActivationData(result) {
		userIDLabel.Text = "ID: available"
	}

	uploadLabel := canvas.NewText("↑ Upload: 0 B/s", components.ColorText())
	uploadLabel.TextSize = components.TextBody
	downloadLabel := canvas.NewText("↓ Download: 0 B/s", components.ColorText())
	downloadLabel.TextSize = components.TextBody

	connectButton := components.NewPrimaryButton("CONNECT", nil)
	connectButtonWrap := container.NewStack(
		components.NewHSpacer(220),
		connectButton,
	)
	isConnected := false
	connectButton.OnTapped = func() {
		isConnected = !isConnected
		if isConnected {
			connectButton.SetText("DISCONNECT")
			statusLabel.Text = "Status: Connected"
			statusLabel.Color = components.ColorStatusConnected()
			statusLabel.Refresh()
			return
		}
		connectButton.SetText("CONNECT")
		statusLabel.Text = "Status: Disconnected"
		statusLabel.Color = components.ColorStatusDisconnected()
		statusLabel.Refresh()
	}

	resetKeyButton := components.NewSecondaryButton("RESET KEY", func() {
		dialog.ShowInformation("Reset key", "Сброс ключа будет доступен в следующей версии.", window)
	})

	content := container.NewVBox(
		titleLabel,
		components.NewVSpacer(components.Spacing12),
		statusLabel,
		components.NewVSpacer(components.Spacing8),
		userIDLabel,
		components.NewVSpacer(components.Spacing16),
		uploadLabel,
		components.NewVSpacer(components.Spacing8),
		downloadLabel,
		components.NewVSpacer(components.Spacing20),
		connectButtonWrap,
		components.NewVSpacer(components.Spacing8),
		resetKeyButton,
	)

	window.SetContent(
		container.NewBorder(
			components.NewVSpacer(components.Spacing16),
			nil,
			components.NewHSpacer(components.Spacing16),
			nil,
			content,
		),
	)
}

func hasAnyActivationData(result core.ActivateResult) bool {
	return strings.TrimSpace(result.SessionToken) != "" ||
		strings.TrimSpace(result.VPNProfile) != "" ||
		strings.TrimSpace(result.ProfileURL) != ""
}

func showSettingsScreen(window fyne.Window, result core.ActivateResult, appSettings *settings.Settings) {
	titleLabel := canvas.NewText("Settings", components.ColorText())
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.TextSize = components.TextTitle

	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		showMainScreen(window, result, appSettings)
	})
	backButton.Importance = widget.LowImportance

	autoConnectToggle := components.NewToggleSwitch(appSettings.Connection.AutoConnectOnLaunch, func(checked bool) {
		appSettings.Connection.AutoConnectOnLaunch = checked
		_ = settings.Save(*appSettings)
	})
	autoReconnectToggle := components.NewToggleSwitch(appSettings.Connection.AutoReconnect, func(checked bool) {
		appSettings.Connection.AutoReconnect = checked
		_ = settings.Save(*appSettings)
	})

	intervalSelector := components.NewSegmentedControl(
		[]components.SegmentOption{
			{ID: "5", Label: "5s"},
			{ID: "10", Label: "10s"},
			{ID: "30", Label: "30s"},
		},
		"10",
		func(value string) {
			switch value {
			case "5":
				appSettings.Connection.ReconnectIntervalSecs = 5
			case "30":
				appSettings.Connection.ReconnectIntervalSecs = 30
			default:
				appSettings.Connection.ReconnectIntervalSecs = 10
			}
			_ = settings.Save(*appSettings)
		},
	)
	switch appSettings.Connection.ReconnectIntervalSecs {
	case 5:
		intervalSelector.SetSelected("5")
	case 30:
		intervalSelector.SetSelected("30")
	default:
		intervalSelector.SetSelected("10")
	}

	modeSelector := components.NewSegmentedControl(
		[]components.SegmentOption{
			{ID: string(settings.ConnectionModeAuto), Label: "Auto"},
			{ID: string(settings.ConnectionModeVLESSRealityOnly), Label: "Reality"},
		},
		string(settings.ConnectionModeAuto),
		func(value string) {
			switch settings.ConnectionMode(value) {
			case settings.ConnectionModeVLESSRealityOnly:
				appSettings.Connection.Mode = settings.ConnectionModeVLESSRealityOnly
			default:
				appSettings.Connection.Mode = settings.ConnectionModeAuto
			}
			_ = settings.Save(*appSettings)
		},
	)
	switch appSettings.Connection.Mode {
	case settings.ConnectionModeVLESSRealityOnly:
		modeSelector.SetSelected(string(settings.ConnectionModeVLESSRealityOnly))
	default:
		modeSelector.SetSelected(string(settings.ConnectionModeAuto))
	}

	startWithWindowsToggle := components.NewToggleSwitch(appSettings.App.StartWithWindows, func(checked bool) {
		// Stub only: OS autostart integration is implemented later.
		appSettings.App.StartWithWindows = checked
		_ = settings.Save(*appSettings)
	})
	startWithWindowsToggle.Disable()

	languageSelector := components.NewSegmentedControl(
		[]components.SegmentOption{
			{ID: string(settings.LanguageRU), Label: "RU"},
			{ID: string(settings.LanguageEN), Label: "EN"},
		},
		string(appSettings.App.Language),
		func(value string) {
			switch settings.Language(value) {
			case settings.LanguageEN:
				appSettings.App.Language = settings.LanguageEN
			default:
				appSettings.App.Language = settings.LanguageRU
			}
			_ = settings.Save(*appSettings)
		},
	)
	switch appSettings.App.Language {
	case settings.LanguageEN:
		languageSelector.SetSelected(string(settings.LanguageEN))
	default:
		languageSelector.SetSelected(string(settings.LanguageRU))
	}

	rememberDeviceToggle := components.NewToggleSwitch(appSettings.Privacy.RememberDevice, func(checked bool) {
		appSettings.Privacy.RememberDevice = checked
		_ = settings.Save(*appSettings)
	})

	clearDataButton := components.NewDangerSecondaryButton("Clear local data", func() {
		dialog.NewConfirm(
			"Clear local data",
			"Удалить локальные настройки и сбросить параметры по умолчанию?",
			func(confirm bool) {
				if !confirm {
					return
				}
				newDefaults, err := settings.Clear()
				if err != nil {
					dialog.ShowInformation("Ошибка", "Не удалось выполнить очистку данных.", window)
					return
				}

				*appSettings = newDefaults
				autoConnectToggle.SetOn(appSettings.Connection.AutoConnectOnLaunch)
				autoReconnectToggle.SetOn(appSettings.Connection.AutoReconnect)
				rememberDeviceToggle.SetOn(appSettings.Privacy.RememberDevice)
				startWithWindowsToggle.SetOn(appSettings.App.StartWithWindows)

				switch appSettings.Connection.ReconnectIntervalSecs {
				case 5:
					intervalSelector.SetSelected("5")
				case 30:
					intervalSelector.SetSelected("30")
				default:
					intervalSelector.SetSelected("10")
				}
				switch appSettings.Connection.Mode {
				case settings.ConnectionModeVLESSRealityOnly:
					modeSelector.SetSelected(string(settings.ConnectionModeVLESSRealityOnly))
				default:
					modeSelector.SetSelected(string(settings.ConnectionModeAuto))
				}
				switch appSettings.App.Language {
				case settings.LanguageEN:
					languageSelector.SetSelected(string(settings.LanguageEN))
				default:
					languageSelector.SetSelected(string(settings.LanguageRU))
				}
			},
			window,
		).Show()
	})

	connectionSection := makeSettingsCard(
		"Connection",
		components.NewSettingRow("Auto-connect on launch", "", autoConnectToggle),
		components.NewSettingRow("Auto-reconnect", "", autoReconnectToggle),
		components.NewSettingRow("Reconnect interval", "", intervalSelector),
		components.NewSettingRow("Connection mode", "", modeSelector),
	)

	appSection := makeSettingsCard(
		"App",
		components.NewSettingRow("Language", "", languageSelector),
		components.NewSettingRow("Start with Windows", "Будет доступно позже.", startWithWindowsToggle),
	)

	privacySection := makeSettingsCard(
		"Privacy & Security",
		components.NewSettingRow("Remember this device", "", rememberDeviceToggle),
		components.NewSettingRow("Clear local data", "Удаляет локальные настройки.", clearDataButton),
	)

	sections := container.NewVBox(
		connectionSection,
		components.NewVSpacer(components.Spacing16),
		appSection,
		components.NewVSpacer(components.Spacing16),
		privacySection,
	)

	scroll := container.NewVScroll(
		container.NewBorder(
			components.NewVSpacer(components.Spacing16),
			components.NewVSpacer(components.Spacing16),
			components.NewHSpacer(components.Spacing16),
			components.NewHSpacer(components.Spacing16),
			sections,
		),
	)
	header := components.NewHeaderBar(
		container.NewStack(
			components.NewHSpacer(components.Spacing40),
			container.NewCenter(backButton),
		),
		titleLabel,
		components.NewHSpacer(components.Spacing40),
	)

	window.SetContent(container.NewBorder(header, nil, nil, nil, scroll))
}

func makeSettingsCard(title string, rows ...fyne.CanvasObject) *fyne.Container {
	cardTitle := components.NewCardTitle(title)
	objects := make([]fyne.CanvasObject, 0, len(rows)*2+1)
	objects = append(objects, cardTitle, components.NewVSpacer(components.Spacing12))
	for _, row := range rows {
		objects = append(objects, row)
	}

	return components.NewCardWithPadding(
		container.NewVBox(objects...),
		components.Spacing20,
		components.Spacing20,
	)
}

func isDevEnvironment() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("VOLTA_ENV"))) {
	case "dev", "development", "local", "debug":
		return true
	default:
		return false
	}
}

