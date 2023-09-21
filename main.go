package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/haste/cache"
)

type Data struct {
	URL  string `json:"Url"`
	Body interface{}
}

func main() {

	app := app.New()

	window := app.NewWindow("Haste")

	shortcuts(&window)
	window.SetMaster()

	url := &widget.Entry{}
	// TODO -> This needs to be text editor; I'll have to make my own
	// text/code editor, or find a library that does it.
	// TODO -> Cache the body, url and method text

	body := widget.NewMultiLineEntry()

	responseUi := widget.NewLabel("Response")

	jsonTmpData := cache.ReadCacheData()

	bodyData, err := json.MarshalIndent(jsonTmpData, "", " ")
	body.SetText(string(bodyData))

	httpMethod := &widget.Select{Options: []string{"GET", "POST", "PUT"}}

	if err != nil {
		fmt.Println(err)
	}

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

	tab := container.NewTabItem("Body", body)

	appTabs := container.NewAppTabs(tab)

	border := container.NewBorder(header, nil, nil, nil, container.NewHSplit(appTabs, container.NewScroll(responseUi)))

	window.SetContent(border)

	window.Resize(fyne.NewSize(1000, 600))

	window.SetOnClosed(func() {
		cache.WriteToCache(body.Text)
	})
	window.ShowAndRun()
}
func shortcuts(window *fyne.Window) {

	w := *window

	ctrlW := &desktop.CustomShortcut{KeyName: fyne.KeyW, Modifier: fyne.KeyModifierControl}

	w.Canvas().AddShortcut(ctrlW, func(shortcut fyne.Shortcut) {
		w.Close()
	})

}
