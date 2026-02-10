package gui

import (
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

	accessInputEntry := widget.NewEntry()
	accessInputEntry.SetPlaceHolder("https://vvpn.io/...")

	statusLabel := widget.NewLabel("")

	continueButton := widget.NewButton("Продолжить", func() {
		statusMessage, ok := core.ValidateAccessInput(accessInputEntry.Text)
		statusLabel.SetText(statusMessage)

		if ok {
			// TODO: auth + VPN
		}
	})

	form := container.NewVBox(
		widget.NewLabel("Ключ доступа или ссылка"),
		accessInputEntry,
		statusLabel,
		continueButton,
		layout.NewSpacer(),
	)

	formContainer := container.New(
		layout.NewVBoxLayout(),
		container.NewPadded(form),
	)

	window.SetContent(formContainer)
	window.Resize(fyne.NewSize(420, 220))
	window.CenterOnScreen()

	window.ShowAndRun()
}
