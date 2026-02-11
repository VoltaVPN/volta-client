package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func NewCardWithPadding(content fyne.CanvasObject, horizontal, vertical float32) *fyne.Container {
	padded := container.NewBorder(
		NewVSpacer(vertical), NewVSpacer(vertical),
		NewHSpacer(horizontal), NewHSpacer(horizontal),
		content,
	)

	// Flat layout: no card shadows/rounded layers.
	return padded
}

func NewHSpacer(width float32) fyne.CanvasObject {
	r := canvas.NewRectangle(color.Transparent)
	r.SetMinSize(fyne.NewSize(width, 1))
	return r
}

func NewVSpacer(height float32) fyne.CanvasObject {
	r := canvas.NewRectangle(color.Transparent)
	r.SetMinSize(fyne.NewSize(1, height))
	return r
}
