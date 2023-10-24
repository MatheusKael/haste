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

	window.SetMainMenu(makeMenu(app, window))

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

		fmt.Println(string(body))

		json.Indent(&dst, body, "", "\t")

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
func makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	newItem := fyne.NewMenuItem("New", nil)
	checkedItem := fyne.NewMenuItem("Checked", nil)
	checkedItem.Checked = true
	disabledItem := fyne.NewMenuItem("Disabled", nil)
	disabledItem.Disabled = true
	otherItem := fyne.NewMenuItem("Other", nil)
	mailItem := fyne.NewMenuItem("Mail", func() { fmt.Println("Menu New->Other->Mail") })
	otherItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("Project", func() { fmt.Println("Menu New->Other->Project") }),
		mailItem,
	)
	fileItem := fyne.NewMenuItem("File", func() { fmt.Println("Menu New->File") })
	dirItem := fyne.NewMenuItem("Directory", func() { fmt.Println("Menu New->Directory") })
	newItem.ChildMenu = fyne.NewMenu("",
		fileItem,
		dirItem,
		otherItem,
	)

	openSettings := func() {
		w := a.NewWindow("Fyne Settings")
		w.Resize(fyne.NewSize(480, 480))
		w.Show()
	}
	settingsItem := fyne.NewMenuItem("Settings", openSettings)
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut
	w.Canvas().AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
		openSettings()
	})

	performFind := func() { fmt.Println("Menu Find") }
	findItem := fyne.NewMenuItem("Find", performFind)
	findItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierShortcutDefault | fyne.KeyModifierAlt | fyne.KeyModifierShift | fyne.KeyModifierControl | fyne.KeyModifierSuper}
	w.Canvas().AddShortcut(findItem.Shortcut, func(shortcut fyne.Shortcut) {
		performFind()
	})

	// a quit item will be appended to our first (File) menu
	file := fyne.NewMenu("File", newItem, checkedItem, disabledItem)
	device := fyne.CurrentDevice()
	if !device.IsMobile() && !device.IsBrowser() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	mainMenu := fyne.NewMainMenu(
		file,
		fyne.NewMenu("Edit", fyne.NewMenuItemSeparator(), findItem),
	)
	checkedItem.Action = func() {
		checkedItem.Checked = !checkedItem.Checked
		mainMenu.Refresh()
	}
	return mainMenu
}

func makeNav() fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return []string{"Children"}
		},
		IsBranch: func(uid string) bool {

			return true
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		//		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
		//			t, ok := tutorials.Tutorials[uid]
		//			if !ok {
		//				fyne.LogError("Missing tutorial panel: "+uid, nil)
		//				return
		//			}
		//			obj.(*widget.Label).SetText(t.Title)
		//			if unsupportedTutorial(t) {
		//				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
		//			} else {
		//				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
		//			}
		//		},
		OnSelected: func(uid string) {
			fmt.Println("Selected: " + uid)
		},
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}
