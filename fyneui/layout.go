package fyneui

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
	"github.com/haste/database"
)

func Build() {

	app := app.New()

	db, err := database.ConnectAndCreateTable()

	defer db.Close()

	window := app.NewWindow("Haste")

	shortcuts(&window)
	window.SetMaster()

	url_field := &widget.Entry{}
	// TODO -> This needs to be text editor; I'll have to make my own
	// text/code editor, or find a library that does it.

	body := widget.NewMultiLineEntry()

	responseUi := widget.NewLabel("Response")

	var jsonTmpData string

	rows, err := database.ReadDataById(1, db)

	httpMethod := &widget.Select{Options: []string{"GET", "POST", "PUT"}}

	for rows.Next() {
		var id int
		var url string
		var method string
		var body string
		var body_format string
		var created_at string
		var updated_at string

		err = rows.Scan(&id, &url, &method, &body, &body_format, &created_at, &updated_at)

		if err != nil {
			fmt.Println(err)
			return
		}

		httpMethod.SetSelected(method)

		url_field.SetText(url)

		jsonTmpData = body
	}

	body.SetText(string(jsonTmpData))

	if err != nil {
		fmt.Println(err)
	}

	button := widget.NewButton("Enter", func() {
		fmt.Println("Click!")

		client := &http.Client{}

		reader := strings.NewReader(body.Text)

		req, err := http.NewRequest(httpMethod.Selected, url_field.Text, reader)

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

	header := container.NewGridWithColumns(3, httpMethod, url_field, container.NewHBox(button))

	tab := container.NewTabItem("Body", body)

	appTabs := container.NewAppTabs(tab)

	border := container.NewBorder(header, nil, nil, nil, container.NewHSplit(appTabs, container.NewScroll(responseUi)))

	window.SetContent(border)

	window.Resize(fyne.NewSize(1000, 600))

	window.SetOnClosed(func() {
		//cache.WriteToCache(body.Text)

		data := &database.RequestData{Body_format: "JSON", Body: body.Text, Method: httpMethod.Selected, Url: "teste.com"}

		rows, err := database.ReadData(db)

		if err != nil {
			fmt.Println(err)
			return
		}

		count := 0

		var id int

		for rows.Next() {
			count++
		}

		if count > 0 {
			fmt.Println(id)
			res, err := database.UpdateData(id, data, db)

			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			rowsAff, err := res.RowsAffected()

			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

			fmt.Println("Rows Affected: ", rowsAff)

		}

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
