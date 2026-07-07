package core

import "gocv.io/x/gocv"

type Transition interface {
	Set(v interface{}) error
	Next(Video, Video, *gocv.Mat) error
}
