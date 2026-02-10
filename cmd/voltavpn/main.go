package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// main is intentionally minimal: no networking, no VPN logic.
// It only wires up the GUI framework and shows a placeholder window.
func main() {
	application := app.New()
	window := application.NewWindow("VoltaVPN")

	window.SetContent(container.NewVBox(
		widget.NewLabel("VoltaVPN client (skeleton)"),
		widget.NewLabel("No VPN logic implemented yet."),
	))

	window.Resize(fyne.NewSize(480, 320))
	window.ShowAndRun()

	// Note: any errors from ShowAndRun are internal to the toolkit.
	// We'll extend error handling when business logic appears.
	if err := recover(); err != nil {
		log.Printf("fatal GUI error: %v", err)
	}
}

