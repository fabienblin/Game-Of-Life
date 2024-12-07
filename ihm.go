package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fabienblin/fynitude"
	"golang.org/x/image/draw"
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
	tappableImage  *fynitude.TappableImageWidget
	tappChan       chan image.Point
}

var _ihm ihm

func initIHM(app fyne.App) fyne.Window {
	window := app.NewWindow("Game Of Life")
	window.SetFullScreen(true)

	_ihm = ihm{
		speedLabel:     binding.NewString(),
		modeLabel:      binding.NewString(),
		menuBackground: canvas.NewRectangle(getModeColor()),
		tappChan:       make(chan image.Point),
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

	rectangle := image.Rect(0, 0, IMAGE_WIDTH, IMAGE_HEIGHT)
	img := image.NewRGBA(rectangle)
	draw.Draw(img, img.Bounds(), &image.Uniform{DEAD}, image.Point{0, 0}, draw.Src)

	imageWidget := fynitude.NewTappableImageWidget(img, nil)
	imageWidget.Tapp = func(pointEvent *fyne.PointEvent) {
		if _ihm.mode != EDIT {
			return
		}
		widgetSize := imageWidget.Size()

		imageBounds := imageWidget.Canvas.Image.Bounds()
		imageWidth := imageBounds.Dx()
		imageHeight := imageBounds.Dy()

		scaleX := float64(widgetSize.Width) / float64(imageWidth)
		scaleY := float64(widgetSize.Height) / float64(imageHeight)

		scale := scaleX
		if scaleY < scaleX {
			scale = scaleY
		}

		renderedWidth := float32(float64(imageWidth) * scale)
		renderedHeight := float32(float64(imageHeight) * scale)

		offsetX := (widgetSize.Width - renderedWidth) / 2
		offsetY := (widgetSize.Height - renderedHeight) / 2

		adjustedX := pointEvent.Position.X - offsetX
		adjustedY := pointEvent.Position.Y - offsetY

		if adjustedX < 0 || adjustedX >= renderedWidth || adjustedY < 0 || adjustedY >= renderedHeight {
			log.Println("Tapped outside the image bounds")
			return
		}

		tappedX := int(float64(adjustedX) / scale)
		tappedY := int(float64(adjustedY) / scale)

		if tappedX < 0 || tappedX >= imageWidth || tappedY < 0 || tappedY >= imageHeight {
			log.Println("Tapped outside the image pixel grid")
			return
		}

		log.Printf("Image tapped at pixel: (%d, %d)", tappedX, tappedY)

		if rgbaImage, ok := imageWidget.Canvas.Image.(*image.RGBA); ok {
			color := rgbaImage.At(tappedX, tappedY)
			log.Printf("Pixel color: %v", color)
		}

		_ihm.tappChan <- image.Point{X: tappedX, Y: tappedY}
	}

	imageWidget.Canvas.FillMode = canvas.ImageFillContain
	imageWidget.Canvas.ScaleMode = canvas.ImageScalePixels

	_ihm.tappableImage = imageWidget

	return container.NewStack(imageWidget)
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

	loadSelectWidget := fynitude.NewDynamicSelectWidget(findGOLimages, triggerLoadImage)

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
