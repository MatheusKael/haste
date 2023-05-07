package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {

	app := app.New()

	window := app.NewWindow("Haste")

	window.SetMaster()

	entry := widget.NewEntry()

	entry.Resize(fyne.NewSize(300, 600))

	button := widget.NewButton("Hello", func() {
		fmt.Println("Hello", entry.Text)
	})

	sendContent := container.NewHBox(entry, button)

	sendContent.Resize(fyne.NewSize(300, 600))

	content := container.NewVBox(sendContent)

	border := container.NewBorder(container.NewVBox(content, widget.NewSeparator()), nil, nil, nil)

	border.Resize(fyne.NewSize(640, 200))

	window.SetContent(border)

	window.Resize(fyne.NewSize(640, 460))
	window.ShowAndRun()
}
