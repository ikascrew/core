package util

import (
	"bufio"
	"image"
	"os"
	"unsafe"

	"gocv.io/x/gocv"
)

func ResizeImage(m gocv.Mat, w, h int) (gocv.Mat, error) {
	//dst := gocv.NewMatWithSize(w,h,gocv.MatTypeCV8U)
	dst := gocv.NewMat()
	gocv.Resize(m, &dst, image.Point{}, 0.5, 0.5, gocv.InterpolationDefault)
	return dst, nil
}

func WriteImage(f string, m gocv.Mat) error {
	gocv.IMWrite(f, m)
	return nil
}

func bstring(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func sbytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func Input() string {
	std := bufio.NewScanner(os.Stdin)
	std.Scan()
	text := std.Text()
	return text
}
