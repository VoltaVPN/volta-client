package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func NewHeaderBar(left, center, right fyne.CanvasObject) *fyne.Container {
	bg := canvas.NewRectangle(ColorSurface())
	bg.SetMinSize(fyne.NewSize(1, HeaderHeight))
	bg.CornerRadius = Radius16

	topGlow := canvas.NewRectangle(ColorElevationLow())
	topGlow.SetMinSize(fyne.NewSize(1, HeaderHeight))
	topGlow.CornerRadius = Radius16

	line := canvas.NewRectangle(ColorBorder())
	line.SetMinSize(fyne.NewSize(1, 1))

	leftSlot := headerSlot(left)
	rightSlot := headerSlot(right)

	row := container.NewBorder(
		nil, nil,
		leftSlot, rightSlot,
		container.NewCenter(center),
	)

	content := container.NewBorder(
		nil, line,
		nil, nil,
		container.NewBorder(NewVSpacer(Spacing16), NewVSpacer(Spacing12), NewHSpacer(Spacing16), NewHSpacer(Spacing16), row),
	)

	return container.NewStack(topGlow, bg, content)
}

func headerSlot(content fyne.CanvasObject) fyne.CanvasObject {
	anchor := canvas.NewRectangle(color.Transparent)
	anchor.SetMinSize(fyne.NewSize(HeaderSideSlotWidth, 1))
	return container.NewStack(anchor, container.NewCenter(content))
}
