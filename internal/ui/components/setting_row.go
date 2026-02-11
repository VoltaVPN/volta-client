package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type settingRowLayout struct{}

func (l *settingRowLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 2 {
		return
	}
	left := objects[0]
	right := objects[1]

	rightSize := right.MinSize()
	leftWidth := size.Width - rightSize.Width - Spacing16
	if leftWidth < 80 {
		leftWidth = 80
	}

	leftSize := fyne.NewSize(leftWidth, size.Height)
	rightY := (size.Height - rightSize.Height) / 2

	left.Move(fyne.NewPos(0, 0))
	left.Resize(leftSize)

	right.Move(fyne.NewPos(size.Width-rightSize.Width, rightY))
	right.Resize(rightSize)
}

func (l *settingRowLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) < 2 {
		return fyne.NewSize(200, 52)
	}
	left := objects[0].MinSize()
	right := objects[1].MinSize()
	height := left.Height
	if right.Height > height {
		height = right.Height
	}
	if height < 52 {
		height = 52
	}
	return fyne.NewSize(left.Width+right.Width+Spacing16, height)
}

func NewSettingRow(title, subtitle string, control fyne.CanvasObject) *fyne.Container {
	titleLabel := canvas.NewText(title, ColorText())
	titleLabel.TextSize = BodyTextSize

	left := fyne.CanvasObject(titleLabel)
	if subtitle != "" {
		subLabel := canvas.NewText(subtitle, ColorTextMuted())
		subLabel.TextSize = CaptionTextSize
		left = container.NewVBox(titleLabel, subLabel)
	}

	rowContent := container.New(&settingRowLayout{}, left, control)

	separator := canvas.NewRectangle(ColorBorder())
	separator.SetMinSize(fyne.NewSize(1, 1))

	return container.NewVBox(
		rowContent,
		NewVSpacer(Spacing8),
		separator,
	)
}

func NewCardTitle(text string) *canvas.Text {
	title := canvas.NewText(text, ColorTextMuted())
	title.TextSize = TextLabel
	title.TextStyle = fyne.TextStyle{Bold: true}
	return title
}
