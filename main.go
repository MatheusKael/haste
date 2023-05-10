package main

import (
	"fmt"
	"image/color"
	"io"
	"net/http"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {

	app := app.New()

	window := app.NewWindow("Haste")

	window.SetMaster()

	url := &widget.Entry{}
	// TODO -> This needs to be text editor,
	// FIX -> Text grid doesn't works like a input, and the documentation
	// sucks
	body := widget.NewTextGrid()

	body.CreateRenderer()

	body.ShowLineNumbers = true

	responseUi := canvas.NewText("Response", color.Black)

	responseUi.Alignment = fyne.TextAlignCenter

	httpMethod := &widget.Select{Options: []string{"GET", "POST", "PUT"}}

	body.SetText("Test")
	button := widget.NewButton("Enter", func() {
		fmt.Println("Click!")

		client := &http.Client{}

		reader := strings.NewReader(body.Text())

		req, err := http.NewRequest(httpMethod.Selected, url.Text, reader)

		if err != nil {
			fmt.Println(err)
		}

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
		}

		body, err := io.ReadAll(resp.Body)

		responseUi.Text = string(body)
	})

	button.Resize(fyne.NewSize(20, 30))

	header := container.NewGridWithColumns(3, httpMethod, url, container.NewHBox(button))

	border := container.NewBorder(header, nil, nil, nil, container.NewHSplit(body, responseUi))

	window.SetContent(border)

	window.Resize(fyne.NewSize(1000, 600))

	window.ShowAndRun()
}
