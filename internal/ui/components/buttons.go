package components

import "fyne.io/fyne/v2/widget"

func NewPrimaryButton(label string, onTap func()) *widget.Button {
	btn := widget.NewButton(label, onTap)
	btn.Importance = widget.HighImportance
	return btn
}

func NewSecondaryButton(label string, onTap func()) *widget.Button {
	btn := widget.NewButton(label, onTap)
	btn.Importance = widget.MediumImportance
	return btn
}

func NewDangerSecondaryButton(label string, onTap func()) *widget.Button {
	// Keep default button style but clearly communicate destructive intent in label.
	return NewSecondaryButton(label, onTap)
}
