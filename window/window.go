package window

import (
	"github.com/ikascrew/core"
	"gocv.io/x/gocv"
)

type Window struct {
	owner *gocv.Window
}

func New(name string) (*Window, error) {
	win := gocv.NewWindow(name)
	w := &Window{
		owner: win,
	}
	return w, nil
}

func (w *Window) Play(v core.Video) error {
	for {
		img, err := v.Next()
		if err != nil {
			return nil
		}
		err = w.Show(img)
		if err != nil {
			return nil
		}
		gocv.WaitKey(1)
	}
	return nil
}

func (w *Window) Show(img *gocv.Mat) error {
	w.owner.IMShow(*img)
	return nil
}

func (w *Window) Wait() {
	gocv.WaitKey(33)
}

func (w *Window) Close() {
	w.owner.Close()
}
