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

// stdin は共有の Scanner。呼び出しごとに新しい Scanner を作ると
// 先読みバッファに残った行が捨てられ、パイプ/リダイレクト入力で
// 2回目以降の Input が空になるため、パッケージで1つだけ持つ
var stdin = bufio.NewScanner(os.Stdin)

func Input() string {
	stdin.Scan()
	return stdin.Text()
}
