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

type ihm struct {
	speedLabel     binding.String
	modeLabel      binding.String
	menuBackground *canvas.Rectangle
	mode           MODE
	speed          int
}

var _ihm ihm

func initIHM(app fyne.App) fyne.Window {
	window := app.NewWindow("Game Of Life")
	window.SetFullScreen(true)

	_ihm = ihm{
		speedLabel:     binding.NewString(),
		modeLabel:      binding.NewString(),
		menuBackground: canvas.NewRectangle(getModeColor()),
	}

	imageContainer := initImageContainer(window)
	menuContainer := initMenuContainer(window)

	rootContainer := container.NewBorder(
		nil,            // Top
		menuContainer,  // Bottom (fixed menu)
		nil,            // Left
		nil,            // Right
		imageContainer, // Center (image taking available space)
	)

	window.SetContent(rootContainer)
	triggerPause()

	return window
}

func initImageContainer(window fyne.Window) *fyne.Container {
	_ = window // ignore unused variable warning

	imageWidget := newTappableImageWidget()

	imageContainer := container.NewStack(imageWidget)

	return imageContainer
}

func initMenuContainer(window fyne.Window) *fyne.Container {
	_ihm.menuBackground.SetMinSize(fyne.NewSize(window.Content().Size().Width, float32(MENU_HEIGHT)))

	menuContainer := container.NewStack(
		_ihm.menuBackground,
		container.NewHBox(
			animationControlContainer(),
			layout.NewSpacer(),
			saveLoadContainer(),
			layout.NewSpacer(),
			statusContainer(),
		),
	)
	menuContainer.Resize(fyne.NewSize(window.Content().Size().Width, float32(MENU_HEIGHT)))

	return menuContainer
}

func animationControlContainer() *fyne.Container {
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

	return buttonsContainer
}

func saveLoadContainer() *fyne.Container {
	saveInputWidget := widget.NewEntry()
	saveInputWidget.SetPlaceHolder("file name")

	saveImageFileButton := widget.NewButton("Save", func() {
		triggerSaveImage(saveInputWidget)
	})

	loadSelectWidget := newDynamicSelect(findGOLimages, triggerLoadImage)

	saveLoadContainer := container.NewGridWithColumns(3,
		saveImageFileButton,
		saveInputWidget,
		loadSelectWidget,
	)

	return saveLoadContainer
}

func statusContainer() *fyne.Container {
	sizeWidget := widget.NewLabel(fmt.Sprintf("Size : %dx%d", IMAGE_WIDTH, IMAGE_HEIGHT))
	speedWidget := widget.NewLabelWithData(_ihm.speedLabel)
	modeWidget := widget.NewLabelWithData(_ihm.modeLabel)

	statusContainer := container.NewHBox(
		sizeWidget,
		speedWidget,
		modeWidget,
	)

	return statusContainer
}

func getModeColor() color.NRGBA {
	if _ihm.mode == RUN {
		return color.NRGBA{R: 20, G: 150, B: 40, A: 255}
	} else if _ihm.mode == EDIT {
		return color.NRGBA{R: 20, G: 40, B: 150, A: 255}
	} else {
		return color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	}
}

func triggerPause() {
	_ihm.mode = EDIT
	_ihm.speed = 0
	updateAppState()
}

func triggerPlay() {
	_ihm.mode = RUN
	_ihm.speed = 1
	updateAppState()
}

func triggerFastForward() {
	_ihm.mode = RUN
	if _ihm.speed == 0 {
		_ihm.speed = 1
	}
	if _ihm.speed < MAX_SPEED {
		_ihm.speed *= 2
	}
	updateAppState()
}

func updateAppState() {
	_ihm.modeLabel.Set(_ihm.mode.String())
	_ihm.speedLabel.Set(fmt.Sprintf("Speed : %d", _ihm.speed))

	if _ihm.menuBackground == nil {
		return
	}
	_ihm.menuBackground.FillColor = getModeColor()
	_ihm.menuBackground.Refresh()
}
