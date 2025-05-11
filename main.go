package main

import (
	"flag"
	"fmt"
	imageprocessing "goroutines_pipeline/image_processing"
	"image"
	"os"
	"path/filepath"
	"strings"
)

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
}

type Status struct {
	Success bool
	Path    string
}

func loadImage(paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		// For each input path create a job and add it to
		// the out channel
		for _, p := range paths {
			job := Job{
				// Fixed: Now keeping the original subdirectories structure
				InputPath: p,
				OutPath:   filepath.Join("images", "output", filepath.Base(p)),
			}
			job.Image = imageprocessing.ReadImage(p)
			out <- job
		}
		close(out)
	}()
	return out
}

func resize(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		// For each input job, create a new job after resize and add it to
		// the out channel
		for job := range input { // Read from the channel
			job.Image = imageprocessing.Resize(job.Image, 0.5)
			out <- job
		}
		close(out)
	}()
	return out
}

func convertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input { // Read from the channel
			job.Image = imageprocessing.Grayscale(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func saveImage(input <-chan Job) <-chan Status {
	out := make(chan Status)
	go func() {
		for job := range input {
			err := imageprocessing.WriteImage(job.OutPath, job.Image)

			stat := Status{
				Success: err == nil,
				Path:    job.InputPath,
			}
			out <- stat
		}
		close(out)
	}()
	return out
}

func adjustAlpha(input <-chan Job, factor float64) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			job.Image = imageprocessing.AdjustAlpha(job.Image, factor)
			out <- job
		}
		close(out)
	}()
	return out
}

func increaseBrightness(input <-chan Job, delta int) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input {
			job.Image = imageprocessing.IncreaseBrightness(job.Image, delta)
			out <- job
		}
		close(out)
	}()
	return out
}

// getImagePaths returns a slice of image paths from the given directory
func getImagePaths(dir string) ([]string, error) {
	var paths []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//  skip the output directory
		if info.IsDir() && strings.Contains(path, "images/output") {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			switch ext {
			case ".jpg", ".jpeg", ".png", ".webp":
				paths = append(paths, path)
			}
		}
		return nil
	})

	return paths, err
}

func runPipelineSequentially(paths []string) {
	for _, p := range paths {
		img := imageprocessing.ReadImage(p)
		if img == nil {
			fmt.Println("Skipping image due to error:", p)
			continue
		}

		outDir := "images/output"
		base := strings.TrimSuffix(filepath.Base(p), filepath.Ext(p)) // e.g., "image1"
		outPath := filepath.Join(outDir, base+"_final.png")

		img = imageprocessing.Resize(img, 0.5)
		img = imageprocessing.Grayscale(img)
		img = imageprocessing.IncreaseBrightness(img, 40)
		img = imageprocessing.AdjustAlpha(img, 0.5)

		err := imageprocessing.WriteImage(outPath, img, "png")
		if err == nil {
			fmt.Println("Sequential Success:", outPath)
		} else {
			fmt.Println("Sequential Failed:", outPath)
		}
	}
}

func runPipelineWithGoroutines(paths []string) {
	channel1 := loadImage(paths)
	channel2 := resize(channel1)
	channel3 := convertToGrayscale(channel2)
	channel4 := increaseBrightness(channel3, 40)
	channel5 := adjustAlpha(channel4, 0.5)
	writeResults := saveImage(channel5)

	for result := range writeResults {
		if result.Success {
			fmt.Println("Goroutine Success!", result.Path)
		} else {
			fmt.Println("Goroutine Failed!", result.Path)
		}
	}
}

func runMatrixTests(paths []string) {
	outDir := "images/output"
	os.MkdirAll(outDir, 0755)

	for _, p := range paths {
		img := imageprocessing.ReadImage(p)
		if img == nil {
			fmt.Println("Skipping:", p)
			continue
		}

		matrix := imageprocessing.ImageToRGBAMatrix(img)

		// Gaussian blur test
		matrix.GaussianBlur(7, 2.0)
		blurred := matrix.RGBAMatrixToImage(matrix)
		blurPath := filepath.Join(outDir, strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))+"_blurred.jpg")
		imageprocessing.WriteImage(blurPath, blurred, "jpeg")
		fmt.Println("Saved:", blurPath)
	}
}

func main() {
	// Create the pipeline
	imagePaths, _ := getImagePaths("images")

	// Weather to use concurrency or not
	useConcurrency := flag.Bool("concurrent", true, "Use goroutines")
	flag.Parse()
	if *useConcurrency {
		// using goroutine
		runPipelineWithGoroutines(imagePaths)
	} else {
		// not using goroutine
		runPipelineSequentially(imagePaths)
	}

	runMatrixTests(imagePaths)
}
