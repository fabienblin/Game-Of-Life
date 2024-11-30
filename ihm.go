package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _speedLabel binding.String = binding.NewString()
var _modeLabel binding.String = binding.NewString()
var _menuBackground *canvas.Rectangle = canvas.NewRectangle(getModeColor())

func initIHM(app fyne.App) fyne.Window {
	window := app.NewWindow("Game Of Life")
	window.SetFullScreen(true)

	imageContainer := initImageContainer(window)

	rootContainer := container.NewBorder(
		nil,                       // Top
		initMenuContainer(window), // Bottom (fixed menu)
		nil,                       // Left
		nil,                       // Right
		imageContainer,            // Center (image taking available space)
	)

	window.SetContent(rootContainer)

	return window
}

func initImageContainer(window fyne.Window) *fyne.Container {
	_ = window // ignore unused variable warning

	imageWidget := newTappableImageWidget()

	imageContainer := container.NewStack(imageWidget)

	return imageContainer
}

func initMenuContainer(window fyne.Window) *fyne.Container {
	pauseButton := widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {
		triggerPause()
	})
	playButton := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		triggerPlay()
	})
	fastForwardButton := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {
		triggerFastForward()
	})

	buttonsContainer := container.NewHBox(
		pauseButton,
		playButton,
		fastForwardButton,
	)

	sizeWidget := widget.NewLabel(fmt.Sprintf("Size : %dx%d", IMAGE_WIDTH, IMAGE_HEIGHT))
	updateAppState()
	speedWidget := widget.NewLabelWithData(_speedLabel)
	modeWidget := widget.NewLabelWithData(_modeLabel)

	textsContainer := container.NewHBox(
		sizeWidget,
		speedWidget,
		modeWidget,
	)

	_menuBackground.SetMinSize(fyne.NewSize(window.Content().Size().Width, float32(MENU_HEIGHT)))

	menuContainer := container.NewStack(
		_menuBackground,
		container.NewHBox(
			buttonsContainer,
			layout.NewSpacer(),
			textsContainer,
		),
	)
	menuContainer.Resize(fyne.NewSize(window.Content().Size().Width, float32(MENU_HEIGHT)))

	return menuContainer
}

func getModeColor() color.NRGBA {
	if _mode == RUN {
		return color.NRGBA{R: 20, G: 150, B: 40, A: 255}
	} else if _mode == EDIT {
		return color.NRGBA{R: 20, G: 40, B: 150, A: 255}
	} else {
		return color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	}
}

func triggerPause() {
	_mode = EDIT
	_speed = 0
	updateAppState()
}

func triggerPlay() {
	_mode = RUN
	_speed = 1
	updateAppState()
}

func triggerFastForward() {
	_mode = RUN
	if _speed == 0 {
		_speed = 1
	}
	if _speed < MAX_SPEED {
		_speed *= 2
	}
	updateAppState()
}

func updateAppState() {
	_modeLabel.Set(_mode.String())
	_speedLabel.Set(fmt.Sprintf("Speed : %d", _speed))

	if _menuBackground == nil {
		return
	}
	_menuBackground.FillColor = getModeColor()
	_menuBackground.Refresh()
}
