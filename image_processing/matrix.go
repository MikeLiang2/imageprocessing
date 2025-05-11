package imageprocessing

import (
	"image"
	"image/color"
	"math"
)

// RGBAPixel represents a pixel in RGBA format.
type RGBAPixel struct {
	R, G, B, A float64
}

// RGBAPixel represents a pixel in RGBA format.
type RGBAMatrix struct {
	Data   [][]RGBAPixel
	Width  int
	Height int
}

// ImageToRGBAMatrix converts an image to a 2D slice of RGBAPixel.
func ImageToRGBAMatrix(img image.Image) RGBAMatrix {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	data := make([][]RGBAPixel, height)
	for y := 0; y < height; y++ {
		row := make([]RGBAPixel, width)
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			row[x] = RGBAPixel{
				R: float64(r) / 257,
				G: float64(g) / 257,
				B: float64(b) / 257,
				A: float64(a) / 257,
			}
		}
		data[y] = row
	}
	return RGBAMatrix{
		Data:   data,
		Width:  width,
		Height: height,
	}
}

// RGBAMatrixToImage converts a 2D slice of RGBAPixel back to an image.
func (m RGBAMatrix) RGBAMatrixToImage(matrix RGBAMatrix) image.Image {
	height := matrix.Height
	width := matrix.Width
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			p := matrix.Data[y][x]
			// Clamp values to [0, 255]
			img.Set(x, y, color.NRGBA{
				R: clamp8(int(p.R)),
				G: clamp8(int(p.G)),
				B: clamp8(int(p.B)),
				A: clamp8(int(p.A)),
			})
		}
	}
	return img
}

// deepCopy creates a deep copy of the RGBAMatrix.
func (m RGBAMatrix) deepCopy() RGBAMatrix {
	copy := make([][]RGBAPixel, m.Height)
	for y := range m.Data {
		row := make([]RGBAPixel, m.Width)
		copy[y] = append(row[:0], m.Data[y]...)
	}
	return RGBAMatrix{
		Data:   copy,
		Width:  m.Width,
		Height: m.Height,
	}
}

// generateGaussianKernel generates a Gaussian kernel of the given size and sigma.
func generateGaussianKernel(size int, sigma float64) [][]float64 {
	kernel := make([][]float64, size)
	half := size / 2
	twoSigmaSq := 2.0 * sigma * sigma

	for y := -half; y <= half; y++ {
		row := make([]float64, size)
		for x := -half; x <= half; x++ {
			exponent := -(float64(x*x + y*y)) / twoSigmaSq
			row[x+half] = math.Exp(exponent)
		}
		kernel[y+half] = row
	}
	return kernel
}

// sumKernel sums all the values in the kernel.
func sumKernel(kernel [][]float64) float64 {
	var sum float64
	for _, row := range kernel {
		for _, val := range row {
			sum += val
		}
	}
	return sum
}

// GaussianBlur applies a Gaussian blur to the RGBAMatrix.
func (m *RGBAMatrix) GaussianBlur(kernelSize int, sigma float64) {
	if kernelSize%2 == 0 || kernelSize < 3 {
		kernelSize = 3 // fallback to safe default
	}
	kernel := generateGaussianKernel(kernelSize, sigma)
	normalize := sumKernel(kernel)

	copyData := m.deepCopy()
	half := kernelSize / 2

	for y := half; y < m.Height-half; y++ {
		for x := half; x < m.Width-half; x++ {
			var sumR, sumG, sumB, sumA float64
			for ky := -half; ky <= half; ky++ {
				for kx := -half; kx <= half; kx++ {
					px := copyData.Data[y+ky][x+kx]
					weight := kernel[ky+half][kx+half]
					sumR += px.R * weight
					sumG += px.G * weight
					sumB += px.B * weight
					sumA += px.A * weight
				}
			}
			m.Data[y][x] = RGBAPixel{
				R: sumR / normalize,
				G: sumG / normalize,
				B: sumB / normalize,
				A: sumA / normalize,
			}
		}
	}
}
