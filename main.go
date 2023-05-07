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

	button := widget.NewButton("Hello", func() {
		fmt.Println("Hello", entry.Text)
	})

	split := container.NewHSplit(button, entry)

	window.SetContent(split)

	window.Resize(fyne.NewSize(640, 460))
	window.ShowAndRun()
}
