package gui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/voltavpn/volta-client/internal/core"
)

func Run() {
	application := app.New()
	window := application.NewWindow("VoltaVPN")

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

	continueButton := widget.NewButton("Продолжить", func() {
		statusMessage, ok := core.ValidateAccessInput(accessInputEntry.Text)
		statusLabel.SetText(statusMessage)

		if ok {
			// TODO: auth + VPN
		}
	})
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
	window.Resize(fyne.NewSize(420, 260))
	window.CenterOnScreen()

	window.ShowAndRun()
}
