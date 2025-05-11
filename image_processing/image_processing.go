package imageprocessing

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/nfnt/resize"
)

// ConvertToNRGBA converts an image.Image to *image.NRGBA.
func ToNRGBA(img image.Image) *image.NRGBA {
	b := img.Bounds()
	dst := image.NewNRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			dst.Set(x, y, img.At(x, y))
		}
	}
	return dst
}

// ReadImage reads an image from the specified path and returns it as an image.Image.
// If the image cannot be decoded, it returns nil and logs the error.
// The function now uses the ToNRGBA function to ensure the image is in NRGBA format.
func ReadImage(path string) image.Image {
	inputFile, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open image: %s, error: %v", path, err)
		return nil
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		log.Printf("Failed to decode image: %s, error: %v", path, err)
		return nil
	}

	// Convert to NRGBA if the image is not already in that format
	return ToNRGBA(img)
}

// WriteImage writes an image to the given path in the specified format.
// If format is empty, it tries to infer the format from the file extension.
func WriteImage(path string, img image.Image, format ...string) error {
	if img == nil {
		return fmt.Errorf("no image data to write: %s", path)
	}

	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer outFile.Close()

	// Determine format (use provided or infer from path)
	var f string
	if len(format) > 0 && format[0] != "" {
		f = strings.ToLower(format[0])
	} else {
		ext := strings.ToLower(filepath.Ext(path))
		f = strings.TrimPrefix(ext, ".") // ".jpg" → "jpg"
	}

	switch f {
	case "jpg", "jpeg":
		err = jpeg.Encode(outFile, img, nil)
	case "png":
		err = png.Encode(outFile, img)
	default:
		err = fmt.Errorf("unsupported image format: %s", f)
	}

	if err != nil {
		return fmt.Errorf("failed to encode image %s: %w", path, err)
	}
	return nil
}

// AdjustAlpha adjusts the alpha channel of an image by a given factor.
// The factor should be between 0.0 (fully transparent) and 1.0 (fully opaque).
// If the factor is outside this range, it will be clamped to [0.0, 1.0].
func AdjustAlpha(img image.Image, alphaFactor float64) image.Image {
	bounds := img.Bounds()
	dst := image.NewNRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			original := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			newAlpha := uint8(float64(original.A) * alphaFactor)
			dst.SetNRGBA(x, y, color.NRGBA{
				R: original.R,
				G: original.G,
				B: original.B,
				A: newAlpha,
			})
		}
	}
	return dst
}

func Grayscale(img image.Image) image.Image {
	if img == nil {
		log.Printf("Cannot convert nil image to grayscale")
		return nil
	}

	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	// Convert each pixel to grayscale
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalPixel := img.At(x, y)
			grayPixel := color.GrayModel.Convert(originalPixel)
			grayImg.Set(x, y, grayPixel)
		}
	}
	return grayImg
}

// IncreaseBrightness increases the brightness of an image by a given delta.
// The delta value can be positive (to increase brightness) or negative (to decrease brightness).
// The function clamps the resulting color values to the range [0, 255].
func IncreaseBrightness(img image.Image, delta int) image.Image {
	bounds := img.Bounds()
	out := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// Increase brightness by delta
			// Note: RGBA values are in the range [0, 65535]
			// We need to divide by 256 to get the actual color value
			out.Set(x, y, color.RGBA{
				R: clamp8(int(r>>8) + delta),
				G: clamp8(int(g>>8) + delta),
				B: clamp8(int(b>>8) + delta),
				A: uint8(a >> 8),
			})
		}
	}
	return out
}

// clamp8 clamps an integer value to the range [0, 255] and returns it as uint8.
func clamp8(v int) uint8 {
	if v > 255 {
		return 255
	}
	if v < 0 {
		return 0
	}
	return uint8(v)
}

func Resize(img image.Image, scale float64) image.Image {
	if img == nil {
		log.Printf("Cannot resize nil image")
		return nil
	}
	if scale <= 0 {
		log.Printf("Invalid scale factor: %.2f", scale)
		return img
	}

	bounds := img.Bounds()
	newWidth := uint(float64(bounds.Dx()) * scale)
	newHeight := uint(float64(bounds.Dy()) * scale)

	return resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
}
