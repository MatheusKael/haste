package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const tempDir = "./tmp/body.json"

func main() {

	app := app.New()

	window := app.NewWindow("Haste")

	window.SetMaster()

	url := &widget.Entry{}
	// TODO -> This needs to be text editor; I'll have to make my own
	// text/code editor
	// TODO -> Cache the body, url and method text

	body := widget.NewMultiLineEntry()

	responseUi := widget.NewLabel("Response")
	// FIX -> if there's no file it panics
	bodyData, err := os.OpenFile(tempDir, os.O_RDWR, 0666)
	stats, err := os.Stat(tempDir)

	defer bodyData.Close()

	data := make([]byte, stats.Size())

	_, err = bodyData.Read(data)
	if err != nil {
		panic(err)
	}

	body.SetText(string(data))

	body.OnChanged = func(text string) {

		bodyData.Truncate(1)
		bodyData.Seek(1, 0)

		_, err = bodyData.Write([]byte(text))

		if err != nil {
			panic(err)
		}

	}
	httpMethod := &widget.Select{Options: []string{"GET", "POST", "PUT"}}

	button := widget.NewButton("Enter", func() {
		fmt.Println("Click!")

		client := &http.Client{}

		reader := strings.NewReader(body.Text)

		req, err := http.NewRequest(httpMethod.Selected, url.Text, reader)

		if err != nil {
			fmt.Println(err)
		}

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
		}

		body, err := io.ReadAll(resp.Body)

		var dst bytes.Buffer

		json.Indent(&dst, body, "=", "\t")

		responseUi.SetText(dst.String())
	})

	button.Resize(fyne.NewSize(20, 30))

	header := container.NewGridWithColumns(3, httpMethod, url, container.NewHBox(button))

	border := container.NewBorder(header, nil, nil, nil, container.NewHSplit(body, container.NewScroll(responseUi)))

	window.SetContent(border)

	window.Resize(fyne.NewSize(1000, 600))

	window.ShowAndRun()

}
