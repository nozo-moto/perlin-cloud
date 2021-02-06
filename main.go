package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"strconv"

	_ "image/png"

	"github.com/nozo-moto/perlin-cloud/perlin"
)

var (
	width    = 1920
	height   = 1080
	filename = "./image.png"
)

func setField() []color.Color {
	p := perlin.New(
		width,
		height,
		rand.Int63(),
	)
	field := make([]color.Color, width*height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			noise := int(p.OctavePerlinNoise(x, y))%128 + 128
			field[y*height+x] = color.RGBA{
				uint8(noise),
				uint8(noise),
				uint8(noise),
				255,
			}
		}
	}
	return field
}

func setArgs() (err error) {
	args := os.Args
	if len(args) == 2 {
		filename = args[1]
	}

	if len(args) == 3 {
		width, err = strconv.Atoi(args[2])
		if err != nil {
			return
		}
		return
	}

	if len(args) == 4 {
		height, err = strconv.Atoi(args[3])
		if err != nil {
			return
		}
	}

	return
}

func main() {
	if err := setArgs(); err != nil {
		panic(err)
	}

	field := setField()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, field[y*height+x])
		}
	}
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
