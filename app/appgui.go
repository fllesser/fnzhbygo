package app

import (
	"fnzhbygo/replace"
	"fnzhbygo/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func RunApp() {
	a := app.New()
	a.Settings().SetTheme(&theme.MyTheme{})
	w := a.NewWindow("Text Replace")
	//w.Resize(fyne.Size{Width: 500, Height: 300})
	w.Resize(fyne.Size{Width: 600, Height: 200})
	componentBox := container.NewVBox(replace.GetButton(), replace.GetProcessBar(), replace.GetGuiLog())
	w.SetContent(componentBox)
	w.ShowAndRun()
}
