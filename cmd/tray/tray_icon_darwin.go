//go:build darwin

package tray

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

func clipboardIcon() []byte {
	const size = 22
	img := image.NewNRGBA(image.Rect(0, 0, size, size))

	black := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	clear := color.NRGBA{}

	fill := func(c color.NRGBA, x1, y1, x2, y2 int) {
		for x := x1; x <= x2; x++ {
			for y := y1; y <= y2; y++ {
				if x >= 0 && x < size && y >= 0 && y < size {
					img.SetNRGBA(x, y, c)
				}
			}
		}
	}

	fill(black, 7, 1, 14, 5)
	fill(clear, 9, 2, 12, 4)
	fill(black, 2, 4, 19, 20)
	fill(clear, 5, 8, 16, 9)
	fill(clear, 5, 11, 16, 12)
	fill(clear, 5, 14, 12, 15)

	var buf bytes.Buffer
	png.Encode(&buf, img) //nolint:errcheck
	return buf.Bytes()
}
