package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// VoltaVPN blue theme — color system (no inline colors in UI).
var (
	// Primary blue
	voltaPrimary       = hexNRGBA(0x3B82F6) // main
	voltaPrimaryHover  = hexNRGBA(0x2563EB)
	voltaPrimaryActive = hexNRGBA(0x1D4ED8)
	voltaAccent        = hexNRGBA(0x60A5FA) // light/accent

	// Neutral
	voltaBgDark    = hexNRGBA(0x0F1419) // almost black
	voltaCardBg    = hexNRGBA(0x1A1F26) // slightly lighter than bg
	voltaText      = hexNRGBA(0xF5F5F5) // white
	voltaTextMuted = hexNRGBA(0x9CA3AF) // secondary gray
	voltaBorder    = hexNRGBA(0x2D3748) // subtle border

	// Status badge
	voltaStatusDisconnected = hexNRGBA(0x6B7280) // gray
	voltaStatusConnecting   = voltaAccent        // blue
	voltaStatusConnected    = hexNRGBA(0x22C55E) // green

	// Shadow
	voltaShadow = hexNRGBA(0x000000)
)

// Status colors for badge (used by UI only via these; no inline hex).
func StatusDisconnectedColor() color.Color { return voltaStatusDisconnected }
func StatusConnectingColor() color.Color   { return voltaStatusConnecting }
func StatusConnectedColor() color.Color    { return voltaStatusConnected }
func AccentColor() color.Color             { return voltaAccent }
func CardBackgroundColor() color.Color     { return voltaCardBg }
func CardBorderColor() color.Color         { return voltaBorder }

func hexNRGBA(hex uint32) color.Color {
	return &color.NRGBA{
		R: uint8(hex >> 16),
		G: uint8(hex >> 8),
		B: uint8(hex),
		A: 0xFF,
	}
}

// Layout constants — consistent spacing/padding (no magic numbers in UI).
const (
	PaddingStandard = 16
	PaddingCard     = 12
	SpacingRow      = 8
	SpacingSection  = 16
	RadiusButton    = 8
	RadiusCard      = 10
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
		return voltaAccent
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
		return RadiusCard
	default:
		return t.base.Size(name)
	}
}
