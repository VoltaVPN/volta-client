package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type ToggleSwitch struct {
	widget.DisableableWidget
	on        bool
	hovered   bool
	onChanged func(bool)

	track *canvas.Rectangle
	knob  *canvas.Circle
}

func NewToggleSwitch(initial bool, onChanged func(bool)) *ToggleSwitch {
	s := &ToggleSwitch{
		on:        initial,
		onChanged: onChanged,
	}
	s.ExtendBaseWidget(s)
	return s
}

func (s *ToggleSwitch) On() bool { return s.on }

func (s *ToggleSwitch) SetOn(v bool) {
	if s.on == v {
		return
	}
	s.on = v
	s.Refresh()
}

func (s *ToggleSwitch) Tapped(_ *fyne.PointEvent) {
	if s.Disabled() {
		return
	}
	s.on = !s.on
	if s.onChanged != nil {
		s.onChanged(s.on)
	}
	s.Refresh()
}

func (s *ToggleSwitch) MouseIn(_ *desktop.MouseEvent) {
	s.hovered = true
	s.Refresh()
}

func (s *ToggleSwitch) MouseMoved(_ *desktop.MouseEvent) {}

func (s *ToggleSwitch) MouseOut() {
	s.hovered = false
	s.Refresh()
}

func (s *ToggleSwitch) CreateRenderer() fyne.WidgetRenderer {
	s.track = canvas.NewRectangle(ColorBorder())
	s.track.CornerRadius = Radius16
	s.knob = canvas.NewCircle(color.NRGBA{R: 245, G: 247, B: 250, A: 255})

	return &toggleSwitchRenderer{
		sw:      s,
		objects: []fyne.CanvasObject{s.track, s.knob},
	}
}

type toggleSwitchRenderer struct {
	sw      *ToggleSwitch
	objects []fyne.CanvasObject
}

func (r *toggleSwitchRenderer) Layout(size fyne.Size) {
	h := float32(26)
	if size.Height < h {
		h = size.Height
	}
	w := float32(46)
	if size.Width < w {
		w = size.Width
	}

	trackPos := fyne.NewPos((size.Width-w)/2, (size.Height-h)/2)
	r.sw.track.Move(trackPos)
	r.sw.track.Resize(fyne.NewSize(w, h))

	knobD := h - 4
	r.sw.knob.Resize(fyne.NewSize(knobD, knobD))

	knobX := trackPos.X + 2
	if r.sw.on {
		knobX = trackPos.X + w - knobD - 2
	}
	r.sw.knob.Move(fyne.NewPos(knobX, trackPos.Y+2))
}

func (r *toggleSwitchRenderer) MinSize() fyne.Size {
	return fyne.NewSize(46, 26)
}

func (r *toggleSwitchRenderer) Refresh() {
	if r.sw.Disabled() {
		r.sw.track.FillColor = ColorBorder()
		r.sw.knob.FillColor = color.NRGBA{R: 149, G: 162, B: 183, A: 255}
	} else if r.sw.on {
		if r.sw.hovered {
			r.sw.track.FillColor = ColorPrimaryHover()
		} else {
			r.sw.track.FillColor = ColorPrimary()
		}
		r.sw.knob.FillColor = color.NRGBA{R: 245, G: 247, B: 250, A: 255}
	} else {
		if r.sw.hovered {
			r.sw.track.FillColor = ColorSurfaceLayered()
		} else {
			r.sw.track.FillColor = ColorBorder()
		}
		r.sw.knob.FillColor = color.NRGBA{R: 245, G: 247, B: 250, A: 255}
	}

	r.sw.track.Refresh()
	r.sw.knob.Refresh()
}

func (r *toggleSwitchRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *toggleSwitchRenderer) Destroy()                      {}
