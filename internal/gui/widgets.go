package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// statusBadge shows a status label on a colored rounded background (gray/blue/green).
type statusBadge struct {
	widget.BaseWidget
	text  *canvas.Text
	bg    *canvas.Rectangle
	label string
	clr   color.Color
}

func newStatusBadge(label string, clr color.Color) *statusBadge {
	s := &statusBadge{label: label, clr: clr}
	s.ExtendBaseWidget(s)
	s.text = canvas.NewText(label, color.NRGBA{0xF5, 0xF5, 0xF5, 0xFF})
	s.text.TextSize = 14
	s.text.Alignment = fyne.TextAlignCenter
	s.bg = canvas.NewRectangle(clr)
	s.bg.CornerRadius = RadiusCard
	return s
}

func (s *statusBadge) SetText(label string) {
	if s.label == label {
		return
	}
	s.label = label
	s.text.Text = label
	s.text.Refresh()
}

func (s *statusBadge) SetColor(clr color.Color) {
	s.clr = clr
	s.bg.FillColor = clr
	s.bg.Refresh()
}

func (s *statusBadge) CreateRenderer() fyne.WidgetRenderer {
	return &statusBadgeRenderer{
		badge: s,
		bg:    s.bg,
		text:  s.text,
	}
}

type statusBadgeRenderer struct {
	badge *statusBadge
	bg    *canvas.Rectangle
	text  *canvas.Text
}

func (r *statusBadgeRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.text.Resize(size)
	r.text.Move(fyne.NewPos(0, 0))
}

func (r *statusBadgeRenderer) MinSize() fyne.Size {
	return r.text.MinSize().Add(fyne.NewSize(PaddingCard*2, PaddingCard))
}

func (r *statusBadgeRenderer) Refresh() {
	r.text.Text = r.badge.label
	r.text.Refresh()
	r.bg.FillColor = r.badge.clr
	r.bg.CornerRadius = RadiusCard
	r.bg.Refresh()
}

func (r *statusBadgeRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.text}
}

func (r *statusBadgeRenderer) Destroy() {}
