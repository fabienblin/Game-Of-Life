package main

import (
	"image"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

var _tappableImage *tappableImageWidget

type tappableImageWidget struct {
	widget.BaseWidget

	canvas   *canvas.Image
	tappChan chan image.Point
}

func (w *tappableImageWidget) Tapped(pointEvent *fyne.PointEvent) {
	if _ihm.mode != EDIT {
		return
	}
	widgetSize := w.Size()

	imageBounds := w.canvas.Image.Bounds()
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

	if rgbaImage, ok := w.canvas.Image.(*image.RGBA); ok {
		color := rgbaImage.At(tappedX, tappedY)
		log.Printf("Pixel color: %v", color)
	}

	w.tappChan <- image.Point{X: tappedX, Y: tappedY}
}

func (w *tappableImageWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(w.canvas)
}

func newTappableImageWidget() *tappableImageWidget {
	// rectangle := image.Rect(0, 0, IMAGE_WIDTH, IMAGE_HEIGHT)
	// m := image.NewRGBA(rectangle)
	m := GenerateRandomImage(IMAGE_WIDTH, IMAGE_HEIGHT)

	imageCanvas := canvas.NewImageFromImage(m)
	imageCanvas.FillMode = canvas.ImageFillContain
	imageCanvas.ScaleMode = canvas.ImageScalePixels

	widget := &tappableImageWidget{canvas: imageCanvas, tappChan: make(chan image.Point)}
	widget.ExtendBaseWidget(widget)

	_tappableImage = widget

	return widget
}
