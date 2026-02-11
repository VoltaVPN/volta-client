package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func NewCardWithPadding(content fyne.CanvasObject, horizontal, vertical float32) *fyne.Container {
	shadow := canvas.NewRectangle(ColorElevationMedium())
	shadow.CornerRadius = Radius20

	bg := canvas.NewRectangle(ColorSurfaceLayered())
	bg.CornerRadius = Radius20

	border := canvas.NewRectangle(ColorBorder())
	border.CornerRadius = Radius20

	padded := container.NewBorder(
		NewVSpacer(vertical), NewVSpacer(vertical),
		NewHSpacer(horizontal), NewHSpacer(horizontal),
		content,
	)

	body := container.NewStack(border, container.NewPadded(bg))
	return container.NewStack(
		shadow,
		container.NewBorder(NewVSpacer(Spacing4), nil, NewHSpacer(Spacing4), nil, body),
		padded,
	)
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
