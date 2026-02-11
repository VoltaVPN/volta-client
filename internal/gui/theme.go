package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/voltavpn/volta-client/internal/ui/components"
)

// VoltaVPN blue theme — color system (no inline colors in UI).
var (
	// Primary blue
	voltaPrimary       = asNRGBA(components.ColorPrimary())
	voltaPrimaryHover  = asNRGBA(components.ColorPrimaryHover())
	voltaPrimaryActive = hexNRGBA(0x1D4ED8)
	voltaAccent        = asNRGBA(components.ColorStatusConnecting())

	// Neutral
	voltaBgDark    = asNRGBA(components.ColorBackground())
	voltaCardBg    = asNRGBA(components.ColorSurface())
	voltaText      = asNRGBA(components.ColorText())
	voltaTextMuted = asNRGBA(components.ColorTextMuted())
	voltaBorder    = asNRGBA(components.ColorBorder())

	// Status badge
	voltaStatusDisconnected = asNRGBA(components.ColorStatusDisconnected())
	voltaStatusConnecting   = voltaAccent        // blue
	voltaStatusConnected    = asNRGBA(components.ColorStatusConnected())

	// Shadow
	voltaShadow = hexNRGBA(0x000000)
)

// Status colors and shared palette helpers (UI не должен использовать hex напрямую).
func StatusDisconnectedColor() color.Color { return voltaStatusDisconnected }
func StatusConnectingColor() color.Color   { return voltaStatusConnecting }
func StatusConnectedColor() color.Color    { return voltaStatusConnected }
func AccentColor() color.Color             { return voltaAccent }
func CardBackgroundColor() color.Color     { return voltaCardBg }
func CardBorderColor() color.Color         { return voltaBorder }
func TextColor() color.Color               { return voltaText }
func TextMutedColor() color.Color          { return voltaTextMuted }

func hexNRGBA(hex uint32) color.Color {
	return &color.NRGBA{
		R: uint8(hex >> 16),
		G: uint8(hex >> 8),
		B: uint8(hex),
		A: 0xFF,
	}
}

func asNRGBA(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()
	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}

// Layout constants — consistent spacing/padding (no magic numbers in UI).
const (
	PaddingStandard = components.Spacing16
	PaddingCard     = components.Spacing12
	SpacingRow      = components.Spacing8
	SpacingSection  = components.Spacing16
	RadiusButton    = components.Radius8
	RadiusCard      = components.Radius12
	RadiusSelection = components.Radius8
	ShadowOffset    = 2
)

// VoltaTheme implements fyne.Theme for VoltaVPN blue theme.
type VoltaTheme struct {
	base fyne.Theme
}

func NewVoltaTheme() *VoltaTheme {
	return &VoltaTheme{base: theme.DefaultTheme()}
}

func (t *VoltaTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Dark variant only; we force dark UI.
	switch name {
	case theme.ColorNameBackground:
		return voltaBgDark
	case theme.ColorNameForeground:
		return voltaText
	case theme.ColorNameDisabled:
		return voltaTextMuted
	case theme.ColorNamePrimary:
		return voltaPrimary
	case theme.ColorNameHover:
		return voltaPrimaryHover
	case theme.ColorNamePressed:
		return voltaPrimaryActive
	case theme.ColorNameFocus:
		return voltaAccent
	case theme.ColorNameInputBackground:
		return voltaCardBg
	case theme.ColorNameInputBorder:
		return voltaBorder
	case theme.ColorNameButton:
		return voltaPrimary
	case theme.ColorNameDisabledButton:
		return voltaTextMuted
	case theme.ColorNamePlaceHolder:
		return voltaTextMuted
	case theme.ColorNameForegroundOnPrimary:
		return voltaText
	case theme.ColorNameSelection:
		// Делаем цвет выделения ближе к фону, чтобы круг вокруг чекбоксов
		// при наведении не был ярким и навязчивым.
		return voltaBgDark
	case theme.ColorNameShadow:
		return voltaShadow
	case theme.ColorNameSeparator:
		return voltaBorder
	case theme.ColorNameOverlayBackground:
		return voltaBgDark
	case theme.ColorNameScrollBar:
		return voltaBorder
	default:
		return t.base.Color(name, variant)
	}
}

func (t *VoltaTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.base.Font(style)
}

func (t *VoltaTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.base.Icon(name)
}

func (t *VoltaTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return PaddingStandard
	case theme.SizeNameInnerPadding:
		return PaddingCard
	case theme.SizeNameInputRadius:
		return RadiusCard
	case theme.SizeNameSelectionRadius:
		return RadiusSelection
	default:
		return t.base.Size(name)
	}
}
