package replace

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var (
	GuiLog     *widget.Entry
	ProcessBar *widget.ProgressBar
	MainButton *widget.Button
	LogString  = ""
)

func GetProcessBar() *widget.ProgressBar {
	ProcessBar = widget.NewProgressBar()
	ProcessBar.Min = 0
	return ProcessBar
}

func GetGuiLog() *widget.Entry {
	GuiLog = widget.NewMultiLineEntry()
	GuiLog.SetMinRowsVisible(12)
	GuiLog.Resize(fyne.Size{Width: 600, Height: 0})
	return GuiLog
}

func GetButton() *widget.Button {
	MainButton = widget.NewButton("开始替换", func() {
		// Do some work here
		MainButton.Disable()
		RepAndWrite()
		// Enable the button
	})
	return MainButton
}

func Log(text string) {
	LogString += text + "\n"
	GuiLog.SetText(LogString)
}

func LogErr(text string, err error) {
	LogString += text + ", 错误信息: " + err.Error() + "\n"
	GuiLog.SetText(LogString)
}
