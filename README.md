# Go Image Processing Pipeline with Concurrency

This project is a replication and extension of [Amrit Singh's image processing pipeline](https://github.com/code-heim/go_21_goroutines_pipeline), implemented in Go. It demonstrates how to use goroutines and channels to process multiple image files efficiently, with support for grayscale conversion, brightness adjustment, alpha transparency, Gaussian blur, and resizing.

---

## Features

-  Concurrency with goroutines (channel-based pipeline)
-  Sequential fallback mode (`-concurrent=false`)
-  Resize with aspect ratio support
-  Image enhancement: grayscale, brightness, alpha
-  Extended: matrix operations, Gaussian blur
-  Multi-format support: `.jpg`, `.png`, `.webp`, `.bmp`, `.tiff`, etc.
-  Error handling for image I/O
-  Unit tests and benchmark testing

# Examples output
Output images are saved to images/output/ with suffixes like:

_final.png → after grayscale + brightness + alpha
_blurred.jpg, _rotated.jpg → for matrix tests
## Original Image
![Original Image](images/test.png)

### Gaussian blur
![Processed Image1](images/output/test_blurred.jpg)

### Scale + Transparency Adjustment + brightness adjustment
![Processed Image2](images/output/test_final.png)

---
## How to Run
## Process Images with Goroutines

go run main.go -concurrent=true

## Process Sequentially
go run main.go -concurrent=false

## Run Unit Tests
go test -v ./...

## Run Benchmarks
go test -bench=. -benchtime=5x

## Benchmarks Results：
Using goroutines  
![Result 1](result1.png)  

Not using goroutines  
![Result 2](result2.png)  

The benchmark results clearly demonstrate that the concurrent image processing pipeline using goroutines performs significantly better than the sequential version:

Concurrent pipeline: 495,122,367 ns/op
Sequential pipeline: 1,221,597,600 ns/op

This  the concurrent version is approximately 2.47× faster than the sequential version under the same image set and processing logic (grayscale → brightness enhancement → alpha adjustment). This improvement is due to effective use of Go’s channel-based goroutines, which allow different stages of the pipeline to process images in parallel.

## GenAI Tools
This project used ChatGPT (OpenAI GPT-4) assisting for:
Refactor pipeline logic
Add modular matrix operations
Write unit tests and benchmark code
Editing README
https://chatgpt.com/share/682015c4-4bb8-8008-a7fa-290da25a887d


# Reference:
# Episode #21: Concurrency in Go: Pipeline Pattern

[Episode link](https://www.codeheim.io/courses/Episode-21-Concurrency-in-Go-Pipeline-Pattern-65c3ca14e4b0628a4e002201)

Requires Golang 1.20 or higher.
