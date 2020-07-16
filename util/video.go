package util

import (
	"bytes"
	"fmt"
	"strings"

	"gocv.io/x/gocv"
)

type Video struct {
	owner *gocv.VideoCapture
	image *gocv.Mat

	name   string
	Width  int
	Height int
	FPS    float64
	FOURCC float64
	Frames int
}

func NewVideo(f string) (*Video, error) {

	width := 0
	height := 0
	fps := 0.0
	frames := 0.0
	fourcc := 0.0

	var owner *gocv.VideoCapture
	var image *gocv.Mat

	if isImage(f) {
		img := gocv.IMRead(f, gocv.IMReadColor)

		image = &img
		width = img.Cols()
		height = img.Rows()
		fps = 33.3
		frames = 30
		fourcc = 1.0
	} else {

		cap, err := gocv.VideoCaptureFile(f)
		if err != nil {
			return nil, err
		}

		owner = cap

		width = int(cap.Get(gocv.VideoCaptureFrameWidth))
		height = int(cap.Get(gocv.VideoCaptureFrameHeight))
		fps = cap.Get(gocv.VideoCaptureFPS)
		frames = cap.Get(gocv.VideoCaptureFrameCount)
		fourcc = cap.Get(gocv.VideoCaptureFOURCC)
	}

	v := &Video{
		owner:  owner,
		image:  image,
		name:   f,
		Width:  width,
		Height: height,
		FPS:    fps,
		FOURCC: fourcc,
		Frames: int(frames),
	}

	return v, nil
}

func isImage(f string) bool {

	if strings.Index(f, ".png") != -1 ||
		strings.Index(f, ".jpg") != -1 ||
		strings.Index(f, ".jpeg") != -1 {
		return true
	}

	return false
}

func (v *Video) Close() {

	if v.image != nil {
		err := v.image.Close()
		if err != nil {
			fmt.Println(err)
		}
	}

	if v.owner != nil {
		err := v.owner.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (v *Video) GetImage(f float64) (*gocv.Mat, error) {

	if v.image != nil {
		return v.image, nil
	}

	if !v.owner.IsOpened() {
		return nil, fmt.Errorf("Capture Open Error")
	}

	m := gocv.NewMat()
	v.owner.Set(gocv.VideoCapturePosFrames, f)
	v.owner.Read(&m)
	return &m, nil
}

func (v *Video) String() string {
	w := bytes.NewBufferString("")
	fmt.Fprintf(w, "File   :[%s] {\n", v.name)
	fmt.Fprintf(w, "  Width  :[%d]\n", v.Width)
	fmt.Fprintf(w, "  Height :[%d]\n", v.Height)
	fmt.Fprintf(w, "  FPS    :[%f]\n", v.FPS)
	fmt.Fprintf(w, "  Frames :[%d]\n", v.Frames)
	fmt.Fprintf(w, "  FOURCC :[%f]\n", v.FOURCC)
	fmt.Fprintf(w, "}\n")
	return w.String()
}
