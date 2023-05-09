package main

import (
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {

	app := app.New()

	window := app.NewWindow("Haste")

	window.SetMaster()

	entry := &widget.Entry{}

	button := widget.NewButton("Enter", func() {
		fmt.Println("Click!")

		client := &http.Client{}

		req, err := http.NewRequest("GET", "https://google.com/search?q=teste", nil)

		if err != nil {
			fmt.Println(err)
		}

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
		}

		body, err := io.ReadAll(resp.Body)
		fmt.Println(string(body))

	})

	button.Resize(fyne.NewSize(20, 30))

	header := container.NewGridWithColumns(3, layout.NewSpacer(), entry, container.NewHBox(button))

	border := container.NewBorder(header, nil, nil, nil)

	window.SetContent(border)

	window.Resize(fyne.NewSize(600, 600))

	window.ShowAndRun()
}
