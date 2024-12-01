package main

import (
	"image"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/draw"
)

var _tappableImage *tappableImageWidget

type tappableImageWidget struct {
	widget.BaseWidget

	canvas *canvas.Image
	tappChan chan image.Point
}

func (w *tappableImageWidget) Tapped(pointEvent *fyne.PointEvent) {
	if _mode != EDIT {
		return
	}
	// Get the widget size (image container size)
	widgetSize := w.Size()

	// Get the original image dimensions
	imageBounds := w.canvas.Image.Bounds() // Assumes `w.image` is an `*image.RGBA` or similar
	imageWidth := imageBounds.Dx()
	imageHeight := imageBounds.Dy()

	// Calculate the scale factors for width and height
	scaleX := float64(widgetSize.Width) / float64(imageWidth)
	scaleY := float64(widgetSize.Height) / float64(imageHeight)

	// Determine the scaling mode (e.g., uniform scaling)
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	// Calculate the rendered size of the image
	renderedWidth := float32(float64(imageWidth) * scale)
	renderedHeight := float32(float64(imageHeight) * scale)

	// Calculate offsets (centered image)
	offsetX := (widgetSize.Width - renderedWidth) / 2
	offsetY := (widgetSize.Height - renderedHeight) / 2

	// Adjust tap position to account for the offset
	adjustedX := pointEvent.Position.X - offsetX
	adjustedY := pointEvent.Position.Y - offsetY

	// Ensure the tap is within the bounds of the rendered image
	if adjustedX < 0 || adjustedX >= renderedWidth || adjustedY < 0 || adjustedY >= renderedHeight {
		log.Println("Tapped outside the image bounds")
		return
	}

	// Map the adjusted coordinates to the original image's pixel coordinates
	tappedX := int(float64(adjustedX) / scale)
	tappedY := int(float64(adjustedY) / scale)

	// Ensure coordinates are within bounds of the original image
	if tappedX < 0 || tappedX >= imageWidth || tappedY < 0 || tappedY >= imageHeight {
		log.Println("Tapped outside the image pixel grid")
		return
	}

	log.Printf("Image tapped at pixel: (%d, %d)", tappedX, tappedY)

	// Access the pixel color
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
	rectangle := image.Rect(0, 0, IMAGE_WIDTH, IMAGE_HEIGHT)
	m := image.NewRGBA(rectangle)
	draw.Draw(m, m.Bounds(), &image.Uniform{DEAD}, image.Point{0, 0}, draw.Src)
	
	imageCanvas := canvas.NewImageFromImage(m)
	imageCanvas.FillMode = canvas.ImageFillContain
	imageCanvas.ScaleMode = canvas.ImageScalePixels

	widget := &tappableImageWidget{canvas: imageCanvas, tappChan: make(chan image.Point)}
	widget.ExtendBaseWidget(widget)

	_tappableImage = widget

	return widget
}
