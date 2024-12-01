package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type MODE int

const (
	EDIT MODE = iota
	RUN
)

func (m MODE) String() string {
	switch m {
	case EDIT:
		return "EDITION"
	default:
		return "RUNNING"
	}
}

const (
	IMAGE_RATIO_WIDTH  int = 1
	IMAGE_RATIO_HEIGHT int = 1
	IMAGE_RATIO_SIZE   int = 50
	IMAGE_WIDTH        int = IMAGE_RATIO_WIDTH * IMAGE_RATIO_SIZE
	IMAGE_HEIGHT       int = IMAGE_RATIO_HEIGHT * IMAGE_RATIO_SIZE
	MENU_HEIGHT        int = 30
	MAX_SPEED          int = 2 * 2 * 2 * 2 // must be a power of 2 for triggerFastForward()
)

var (
	ALIVE  color.RGBA = color.RGBA{255, 255, 255, 255}
	DEAD   color.RGBA = color.RGBA{0, 0, 0, 255}
	_mode  MODE
	_speed int
)

func init() {
	triggerPause()
}

func main() {
	app := app.New()
	window := initIHM(app)

	window.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
		if keyEvent.Name == fyne.KeyEscape {
			app.Quit()
		}
	})

	refreshImageRoutine()
	window.ShowAndRun()
}
