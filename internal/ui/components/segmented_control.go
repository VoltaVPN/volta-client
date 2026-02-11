package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type SegmentOption struct {
	ID    string
	Label string
}

type SegmentedControl struct {
	widget.BaseWidget
	options   []SegmentOption
	selected  string
	hoveredID string
	onChanged func(string)

	bg        *canvas.Rectangle
	selection *canvas.Rectangle
	labels    []*canvas.Text
}

func NewSegmentedControl(options []SegmentOption, selected string, onChanged func(string)) *SegmentedControl {
	s := &SegmentedControl{
		options:   options,
		selected:  selected,
		onChanged: onChanged,
	}
	s.ExtendBaseWidget(s)
	return s
}

func (s *SegmentedControl) SetSelected(id string) {
	if s.selected == id {
		return
	}
	s.selected = id
	s.Refresh()
}

func (s *SegmentedControl) Selected() string {
	return s.selected
}

func (s *SegmentedControl) Tapped(ev *fyne.PointEvent) {
	if len(s.options) == 0 {
		return
	}

	width := s.Size().Width
	if width <= 0 {
		return
	}
	segmentWidth := width / float32(len(s.options))
	index := int(ev.Position.X / segmentWidth)
	if index < 0 || index >= len(s.options) {
		return
	}
	next := s.options[index].ID
	if next == s.selected {
		return
	}
	s.selected = next
	if s.onChanged != nil {
		s.onChanged(next)
	}
	s.Refresh()
}

func (s *SegmentedControl) MouseIn(_ *desktop.MouseEvent) {
	s.hoveredID = ""
	s.Refresh()
}

func (s *SegmentedControl) MouseMoved(ev *desktop.MouseEvent) {
	next := s.segmentIDAt(ev.Position)
	if next == s.hoveredID {
		return
	}
	s.hoveredID = next
	s.Refresh()
}

func (s *SegmentedControl) MouseOut() {
	s.hoveredID = ""
	s.Refresh()
}

func (s *SegmentedControl) segmentIDAt(pos fyne.Position) string {
	if len(s.options) == 0 {
		return ""
	}
	width := s.Size().Width
	if width <= 0 {
		return ""
	}
	segmentWidth := width / float32(len(s.options))
	index := int(pos.X / segmentWidth)
	if index < 0 || index >= len(s.options) {
		return ""
	}
	return s.options[index].ID
}

func (s *SegmentedControl) CreateRenderer() fyne.WidgetRenderer {
	s.bg = canvas.NewRectangle(ColorSurface())
	s.bg.CornerRadius = 0

	s.selection = canvas.NewRectangle(ColorPrimary())
	s.selection.CornerRadius = 0

	s.labels = make([]*canvas.Text, 0, len(s.options))
	objects := make([]fyne.CanvasObject, 0, len(s.options)+2)
	objects = append(objects, s.bg, s.selection)
	for _, opt := range s.options {
		t := canvas.NewText(opt.Label, ColorTextMuted())
		t.TextSize = BodyTextSize
		t.Alignment = fyne.TextAlignCenter
		s.labels = append(s.labels, t)
		objects = append(objects, t)
	}

	return &segmentedRenderer{control: s, objects: objects}
}

type segmentedRenderer struct {
	control *SegmentedControl
	objects []fyne.CanvasObject
}

func (r *segmentedRenderer) Layout(size fyne.Size) {
	r.control.bg.Resize(size)

	if len(r.control.options) == 0 {
		r.control.selection.Hide()
		return
	}

	segmentWidth := size.Width / float32(len(r.control.options))
	selectedIndex := 0
	for i, opt := range r.control.options {
		if opt.ID == r.control.selected {
			selectedIndex = i
			break
		}
	}

	inset := float32(2)
	r.control.selection.Move(fyne.NewPos(float32(selectedIndex)*segmentWidth+inset, inset))
	r.control.selection.Resize(fyne.NewSize(segmentWidth-inset*2, size.Height-inset*2))
	r.control.selection.Show()

	for i, label := range r.control.labels {
		x := float32(i) * segmentWidth
		label.Move(fyne.NewPos(x, (size.Height-label.MinSize().Height)/2))
		label.Resize(fyne.NewSize(segmentWidth, label.MinSize().Height))
	}
}

func (r *segmentedRenderer) MinSize() fyne.Size {
	width := float32(len(r.control.options)) * 68
	if width < 136 {
		width = 136
	}
	return fyne.NewSize(width, ControlHeightDefault)
}

func (r *segmentedRenderer) Refresh() {
	for i, opt := range r.control.options {
		if i >= len(r.control.labels) {
			continue
		}
		lbl := r.control.labels[i]
		lbl.Text = opt.Label
		if opt.ID == r.control.selected {
			lbl.Color = ColorText()
		} else if opt.ID == r.control.hoveredID {
			lbl.Color = ColorPrimaryHover()
		} else {
			lbl.Color = ColorTextMuted()
		}
		lbl.Refresh()
	}
	r.control.bg.Refresh()
	r.control.selection.Refresh()
}

func (r *segmentedRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *segmentedRenderer) Destroy()                      {}
