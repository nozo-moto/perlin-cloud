package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"

	_ "image/png"

	"github.com/aquilax/go-perlin"
)

const (
	//	WIDTH  = 1920
	//	HEIGHT = 1080
	WIDTH  = 100
	HEIGHT = 100
)

func setField() []color.Color {
	const (
		alpha       = 1.
		beta        = 1.
		n           = 3
		seed  int64 = 10000
	)
	p := perlin.NewPerlinRandSource(alpha, beta, n, rand.NewSource(seed))
	field := make([]color.Color, WIDTH*HEIGHT)
	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			noise := p.Noise2D(float64(x)/10, float64(y)/10)
			field[y*HEIGHT+x] = color.RGBA{0, 0, 0, uint8(255 * noise)}
		}
	}
	return field
}

func main() {
	field := setField()
	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			img.Set(x, y, field[y*HEIGHT+x])
		}
	}
	f, err := os.Create("./image.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
