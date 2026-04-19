//go:build windows

package tray

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
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

	return encodeICO(img)
}

// encodeICO encodes an NRGBA image as a Windows .ico file (32bpp, single image).
// LoadImage on Windows only accepts .ico format, not PNG.
func encodeICO(img *image.NRGBA) []byte {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	// XOR mask: BGRA pixels, bottom-up row order
	xorMask := make([]byte, w*h*4)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.NRGBAAt(x, y)
			dstRow := h - 1 - y
			idx := (dstRow*w + x) * 4
			xorMask[idx+0] = c.B
			xorMask[idx+1] = c.G
			xorMask[idx+2] = c.R
			xorMask[idx+3] = c.A
		}
	}

	// AND mask: 1 bit per pixel, padded to 32-bit row boundary, all zeros (opaque)
	andRowBytes := (w + 31) / 32 * 4
	andMask := make([]byte, andRowBytes*h)

	dibSize := 40 + len(xorMask) + len(andMask)

	var buf bytes.Buffer

	// ICONDIR header
	binary.Write(&buf, binary.LittleEndian, uint16(0)) // reserved
	binary.Write(&buf, binary.LittleEndian, uint16(1)) // type = icon
	binary.Write(&buf, binary.LittleEndian, uint16(1)) // count = 1

	// ICONDIRENTRY
	buf.WriteByte(byte(w))
	buf.WriteByte(byte(h))
	buf.WriteByte(0) // colorCount
	buf.WriteByte(0) // reserved
	binary.Write(&buf, binary.LittleEndian, uint16(1))          // planes
	binary.Write(&buf, binary.LittleEndian, uint16(32))         // bitCount
	binary.Write(&buf, binary.LittleEndian, uint32(dibSize))    // bytesInRes
	binary.Write(&buf, binary.LittleEndian, uint32(6+16))       // imageOffset

	// BITMAPINFOHEADER
	binary.Write(&buf, binary.LittleEndian, uint32(40))
	binary.Write(&buf, binary.LittleEndian, int32(w))
	binary.Write(&buf, binary.LittleEndian, int32(h*2)) // doubled height for ICO
	binary.Write(&buf, binary.LittleEndian, uint16(1))  // planes
	binary.Write(&buf, binary.LittleEndian, uint16(32)) // bitCount
	binary.Write(&buf, binary.LittleEndian, uint32(0))  // BI_RGB
	binary.Write(&buf, binary.LittleEndian, uint32(0))  // biSizeImage
	binary.Write(&buf, binary.LittleEndian, int32(0))   // biXPelsPerMeter
	binary.Write(&buf, binary.LittleEndian, int32(0))   // biYPelsPerMeter
	binary.Write(&buf, binary.LittleEndian, uint32(0))  // biClrUsed
	binary.Write(&buf, binary.LittleEndian, uint32(0))  // biClrImportant

	buf.Write(xorMask)
	buf.Write(andMask)

	return buf.Bytes()
}
