package core

import (
	"gocv.io/x/gocv"
)

type Video interface {
	Next() (*gocv.Mat, error)
	Wait() float64

	Set(int)
	Current() int

	Source() string

	Release() error
}
