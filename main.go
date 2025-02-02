package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/achhabra2/kqb-json-viewer/icons"
	"github.com/achhabra2/kqb-json-viewer/stats"
)

func main() {
	names := stats.ListStatFiles()
	data := stats.ReadJson(names[0])
	mainIcon := fyne.NewStaticResource("logo.png", icons.Logo)

	a := app.NewWithID("com.kqb-json-viewer.app")
	a.SetIcon(mainIcon)
	appTheme := myTheme{}
	a.Settings().SetTheme(&appTheme)

	w := a.NewWindow("KQB JSON Viewer")

	kqbApp := KQBApp{
		files:        names,
		selectedData: data,
		a:            a,
		w:            w,
	}

	kqbApp.ShowMainWindow()

}
