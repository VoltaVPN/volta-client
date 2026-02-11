package components

import "image/color"

const (
	// Spacing scale.
	Spacing4  float32 = 4
	Spacing8  float32 = 8
	Spacing12 float32 = 12
	Spacing16 float32 = 16
	Spacing20 float32 = 20
	Spacing24 float32 = 24
	Spacing32 float32 = 32
	Spacing40 float32 = 40

	// Radius scale.
	Radius8    float32 = 8
	Radius10   float32 = 10
	Radius12   float32 = 12
	Radius16   float32 = 16
	Radius20   float32 = 20
	RadiusPill float32 = 999

	// Typography scale.
	TextDisplay  float32 = 30
	TextHeadline float32 = 22
	TextTitle    float32 = 18
	TextBody     float32 = 14
	TextLabel    float32 = 13
	TextCaption  float32 = 12

	// Component sizing.
	HeaderHeight         float32 = 72
	ControlHeightDefault float32 = 36
	ControlHeightLarge   float32 = 72
	ConnectButtonMinW    float32 = 220
	HeaderSideSlotWidth  float32 = 96
	LogoMarkSize         float32 = 36
	StatIconSlotSize     float32 = 16

	// Backward-compatible aliases used by current UI code.
	TitleTextSize   float32 = TextDisplay
	BodyTextSize    float32 = TextBody
	CaptionTextSize float32 = TextCaption
)

const (
	ElevationLowAlpha    uint8 = 34
	ElevationMediumAlpha uint8 = 54
	ElevationHighAlpha   uint8 = 72
)

var (
	bgBase       = hexNRGBA(0x0E131A)
	bgElevated   = hexNRGBA(0x171D26)
	bgLayered    = hexNRGBA(0x1C2532)
	borderSubtle = hexNRGBA(0x293345)

	textPrimary = hexNRGBA(0xF5F7FA)
	textMuted   = hexNRGBA(0x95A2B7)

	brandPrimary   = hexNRGBA(0x3B82F6)
	brandPrimaryHi = hexNRGBA(0x2563EB)

	statusSuccess = hexNRGBA(0x22C55E)
	statusMuted   = hexNRGBA(0x6B7280)
	statusPending = hexNRGBA(0x60A5FA)

	danger = hexNRGBA(0xDC2626)
)

func hexNRGBA(hex uint32) color.NRGBA {
	return color.NRGBA{
		R: uint8(hex >> 16),
		G: uint8(hex >> 8),
		B: uint8(hex),
		A: 0xFF,
	}
}

func ColorWithAlpha(c color.Color, alpha uint8) color.NRGBA {
	r, g, b, _ := c.RGBA()
	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: alpha,
	}
}

func ColorBackground() color.Color { return bgBase }
func ColorSurface() color.Color    { return bgElevated }
func ColorSurfaceLayered() color.Color {
	return bgLayered
}
func ColorBorder() color.Color     { return borderSubtle }
func ColorText() color.Color       { return textPrimary }
func ColorTextMuted() color.Color  { return textMuted }
func ColorElevationLow() color.Color {
	return color.NRGBA{R: 0, G: 0, B: 0, A: ElevationLowAlpha}
}
func ColorElevationMedium() color.Color {
	return color.NRGBA{R: 0, G: 0, B: 0, A: ElevationMediumAlpha}
}
func ColorElevationHigh() color.Color {
	return color.NRGBA{R: 0, G: 0, B: 0, A: ElevationHighAlpha}
}

func ColorPrimary() color.Color      { return brandPrimary }
func ColorPrimaryHover() color.Color { return brandPrimaryHi }

func ColorStatusConnected() color.Color    { return statusSuccess }
func ColorStatusDisconnected() color.Color { return statusMuted }
func ColorStatusConnecting() color.Color   { return statusPending }

func ColorDanger() color.Color { return danger }
