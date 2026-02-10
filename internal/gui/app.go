package gui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	application := app.New()
	window := application.NewWindow("VoltaVPN")

	accessKeyEntry := widget.NewPasswordEntry()
	accessKeyEntry.SetPlaceHolder("Enter access key")

	statusLabel := widget.NewLabel("")

	continueButton := widget.NewButton("Continue", func() {
		key := strings.TrimSpace(accessKeyEntry.Text)
		if key == "" {
			statusLabel.SetText("Access key is required.")
		} else {
			statusLabel.SetText("Access key looks valid (local check only).")
		}
	})

	form := container.NewVBox(
		widget.NewLabel("VoltaVPN"),
		layout.NewSpacer(),
		widget.NewLabel("Access Key"),
		accessKeyEntry,
		continueButton,
		statusLabel,
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
