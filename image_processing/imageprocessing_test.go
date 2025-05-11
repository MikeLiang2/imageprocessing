package imageprocessing

import (
	"image"
	"image/color"
	"testing"
)

func TestGrayscale(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})

	gray := Grayscale(img)
	if _, ok := gray.(*image.Gray); !ok {
		t.Errorf("Expected *image.Gray, got %T", gray)
	}

	val := gray.At(0, 0).(color.Gray).Y
	expected := uint8(76) // 0.299 * 255 = ~76
	if val != expected {
		t.Errorf("Expected gray value %d, got %d", expected, val)
	}
}

func TestAdjustAlpha(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.NRGBA{R: 100, G: 150, B: 200, A: 255})

	adjusted := AdjustAlpha(img, 0.5)
	alpha := adjusted.(*image.NRGBA).NRGBAAt(0, 0).A
	if alpha != 127 {
		t.Errorf("Expected alpha 127, got %d", alpha)
	}
}

func TestIncreaseBrightness(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.NRGBA{R: 100, G: 100, B: 100, A: 255})

	bright := IncreaseBrightness(img, 50)
	px := bright.(*image.RGBA).RGBAAt(0, 0)
	if px.R != 150 || px.G != 150 || px.B != 150 {
		t.Errorf("Expected RGB to be 150, got R:%d G:%d B:%d", px.R, px.G, px.B)
	}
}

func TestImageToRGBAMatrixAndBack(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.NRGBA{R: 10, G: 20, B: 30, A: 255})
	img.Set(1, 0, color.NRGBA{R: 40, G: 50, B: 60, A: 255})
	img.Set(0, 1, color.NRGBA{R: 70, G: 80, B: 90, A: 255})
	img.Set(1, 1, color.NRGBA{R: 100, G: 110, B: 120, A: 255})

	matrix := ImageToRGBAMatrix(img)
	reconverted := matrix.RGBAMatrixToImage(matrix)

	// Check pixel at (1, 1)
	r, g, b, a := reconverted.At(1, 1).RGBA()
	expect := color.NRGBA{R: 100, G: 110, B: 120, A: 255}
	if uint8(r>>8) != expect.R || uint8(g>>8) != expect.G || uint8(b>>8) != expect.B || uint8(a>>8) != expect.A {
		t.Errorf("Expected pixel (1,1) to be %v, got R:%d G:%d B:%d A:%d",
			expect, uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
	}

}

func TestResize(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 100, 200))
	resized := Resize(img, 0.5)

	expectedWidth := 50
	expectedHeight := 100

	if resized.Bounds().Dx() != expectedWidth || resized.Bounds().Dy() != expectedHeight {
		t.Errorf("Expected resized image to be %dx%d, got %dx%d",
			expectedWidth, expectedHeight,
			resized.Bounds().Dx(), resized.Bounds().Dy())
	}
}

func TestToNRGBA(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	converted := ToNRGBA(img)

	if converted == nil {
		t.Errorf("Expected *image.NRGBA, got nil")
	}
}

func TestRGBAMatrixConversion(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.NRGBA{R: 10, G: 20, B: 30, A: 255})
	img.Set(1, 0, color.NRGBA{R: 40, G: 50, B: 60, A: 255})
	img.Set(0, 1, color.NRGBA{R: 70, G: 80, B: 90, A: 255})
	img.Set(1, 1, color.NRGBA{R: 100, G: 110, B: 120, A: 255})

	matrix := ImageToRGBAMatrix(img)
	reconverted := matrix.RGBAMatrixToImage(matrix)

	r, g, b, a := reconverted.At(1, 1).RGBA()
	expect := color.NRGBA{R: 100, G: 110, B: 120, A: 255}
	if uint8(r>>8) != expect.R || uint8(g>>8) != expect.G || uint8(b>>8) != expect.B || uint8(a>>8) != expect.A {
		t.Errorf("Expected pixel (1,1) to be %v, got R:%d G:%d B:%d A:%d", expect, uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
	}
}

func TestGaussianBlurDoesNotPanic(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 5, 5))
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			img.Set(x, y, color.NRGBA{R: 100, G: 100, B: 100, A: 255})
		}
	}
	matrix := ImageToRGBAMatrix(img)

	// Run blur
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GaussianBlur panicked with: %v", r)
		}
	}()
	matrix.GaussianBlur(3, 1.0)

	// Basic check that something changed or stayed within expected bounds
	p := matrix.Data[2][2]
	if p.R < 90 || p.R > 110 {
		t.Errorf("Blurred center pixel R out of expected range: %.2f", p.R)
	}
}
