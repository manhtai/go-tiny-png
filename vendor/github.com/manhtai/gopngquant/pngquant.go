/*
Copyright (c) 2016, The go-imagequant author(s)

Permission to use, copy, modify, and/or distribute this software for any purpose
with or without fee is hereby granted, provided that the above copyright notice
and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND ISC DISCLAIMS ALL WARRANTIES WITH REGARD TO
THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS.
IN NO EVENT SHALL ISC BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR
CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA
OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS
ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS
SOFTWARE.
*/

package pngquant

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/manhtai/imagequant"
)

// GoImageToRgba32 convert Go Image to RGBA32 bytes
func GoImageToRgba32(im image.Image) []byte {
	w := im.Bounds().Max.X
	h := im.Bounds().Max.Y
	ret := make([]byte, w*h*4)

	p := 0

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r16, g16, b16, a16 := im.At(x, y).RGBA()

			ret[p+0] = uint8(r16 >> 8)
			ret[p+1] = uint8(g16 >> 8)
			ret[p+2] = uint8(b16 >> 8)
			ret[p+3] = uint8(a16 >> 8)
			p += 4
		}
	}

	return ret
}

// Rgb8PaletteToGoImage convert from RBG8 byte to Go Image
func Rgb8PaletteToGoImage(w, h int, rgb8data []byte, pal color.Palette) image.Image {
	rect := image.Rectangle{
		Max: image.Point{
			X: w,
			Y: h,
		},
	}

	ret := image.NewPaletted(rect, pal)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			ret.SetColorIndex(x, y, rgb8data[y*w+x])
		}
	}

	return ret
}

// CompressPng take a png data source and output a png data dest, with a speed/quality
// tradeoff
func CompressPng(source io.Reader, dest io.Writer, speed int) error {
	// Png decode
	img, err := png.Decode(source)
	if err != nil {
		return fmt.Errorf("png.Decode: %s", err.Error())
	}

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	// Init libimagequant attributes holder
	attr, err := imagequant.NewAttributes()
	if err != nil {
		return fmt.Errorf("NewAttributes: %s", err.Error())
	}
	defer attr.Release()

	// Set speed
	err = attr.SetSpeed(speed)
	if err != nil {
		return fmt.Errorf("SetSpeed: %s", err.Error())
	}

	rgba32data := GoImageToRgba32(img)

	// New image from RGBA32 data (not Go Image)
	iqm, err := imagequant.NewImage(attr, string(rgba32data), width, height, 0)
	if err != nil {
		return fmt.Errorf("NewImage: %s", err.Error())
	}
	defer iqm.Release()

	// Compress image using above attributes
	res, err := iqm.Quantize(attr)
	if err != nil {
		return fmt.Errorf("Quantize: %s", err.Error())
	}
	defer res.Release()

	// Get RBG8 data from compressed image
	rgb8data, err := res.WriteRemappedImage()
	if err != nil {
		return fmt.Errorf("WriteRemappedImage: %s", err.Error())
	}

	// Convert RBG8 to Go Image
	im2 := Rgb8PaletteToGoImage(res.GetImageWidth(), res.GetImageHeight(),
		rgb8data, res.GetPalette())

	// Write out the image
	penc := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	err = penc.Encode(dest, im2)
	if err != nil {
		return fmt.Errorf("png.Encode: %s", err.Error())
	}

	return nil
}

// CompressPngFile compress file instead of io.Reader
func CompressPngFile(sourceFile string, destFile string, speed int) error {
	s, err := os.OpenFile(sourceFile, os.O_RDONLY, 0444)
	if err != nil {
		log.Fatalf("os.OpenFile: %s", err.Error())
	}
	defer s.Close()

	d, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("os.OpenFile: %s", err.Error())
	}
	defer d.Close()

	return CompressPng(s, d, speed)
}
