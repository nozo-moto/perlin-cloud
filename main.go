package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	_ "image/png"
)

const (
	WIDTH  = 1920
	HEIGHT = 1080
)

func setField() []color.Color {
	field := make([]color.Color, WIDTH*HEIGHT)
	for i := 0; i < len(field); i++ {
		field[i] = color.RGBA{0, 0, 0, 255}
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
