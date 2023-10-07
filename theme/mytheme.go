package theme

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

var (
	//go:embed fonts/yezi.ttf
	Yezi []byte
)

type MyTheme struct{}

var _ fyne.Theme = (*MyTheme)(nil)

func (m *MyTheme) Font(fyne.TextStyle) fyne.Resource {
	return &fyne.StaticResource{
		StaticName:    "yezi.ttf",
		StaticContent: Yezi,
	}
}

func (*MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (*MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*MyTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
