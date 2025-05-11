package main

import (
	"testing"
)

// benchmark tests for two different implementations of the image processing pipeline
// one using goroutines and the other using a sequential approach
func BenchmarkRunPipelineWithGoroutines(b *testing.B) {
	imagePaths, _ := getImagePaths("images")
	for i := 0; i < b.N; i++ {
		runPipelineWithGoroutines(imagePaths)
	}
}

func BenchmarkRunPipelineSequentially(b *testing.B) {
	imagePaths, _ := getImagePaths("images")
	for i := 0; i < b.N; i++ {
		runPipelineSequentially(imagePaths)
	}
}
