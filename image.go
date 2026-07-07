package core

import (
	"image"

	"gocv.io/x/gocv"
)

type Image struct {
	mat  *gocv.Mat
	img  image.Image
	lazy bool
}

func NewImage(w, h int, lazy bool) *Image {
	var img Image
	return &img
}
