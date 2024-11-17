package main

import (
	"image"
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
	IMAGE_RATIO_WIDTH int = 16
	IMAGE_RATIO_HEIGHT int = 9
	IMAGE_RATIO_SIZE int = 20
	IMAGE_WIDTH int = IMAGE_RATIO_WIDTH * IMAGE_RATIO_SIZE
	IMAGE_HEIGHT int = IMAGE_RATIO_HEIGHT * IMAGE_RATIO_SIZE
	MENU_HEIGHT int = 30
)

var (
	_mode MODE
	_image *image.RGBA
	_speed int
)

func init() {
	onPause()
	_image = image.NewRGBA(image.Rect(0, 0, IMAGE_WIDTH, IMAGE_HEIGHT))
}

func main() {
	window := initIHM()
	window.ShowAndRun()
}