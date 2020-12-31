package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"

	_ "image/png"
)

const (
	WIDTH  = 1920
	HEIGHT = 1080
	//WIDTH  = 100
	//HEIGHT = 100
)

type PerlinNoise struct {
	hashtable []int
}

// Fade 6t^5 - 15t^4 + 10t^3.
func (p *PerlinNoise) Fade(t float64) float64 {
	return (6 * math.Pow(t, 5)) - (15 * math.Pow(5, 4)) + (10 * math.Pow(t, 3))
}

// Lerp 線形補間
func (p *PerlinNoise) Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func (p *PerlinNoise) setHash(seed int64) {
	p.hashtable = make([]int, WIDTH*HEIGHT)
	randtable := make([]int, HEIGHT)
	for i := 0; i < HEIGHT; i++ {
		randtable[i] = rand.Intn(HEIGHT)
	}
	rand.Seed(seed)
	for i := 0; i < WIDTH*HEIGHT; i++ {
		p.hashtable[i] = randtable[i%HEIGHT]
	}
}

func (p *PerlinNoise) getHash(x, y int) int {
	//	x %= WIDTH
	//	y %= HEIGHT
	x %= 255
	y %= 255
	return p.hashtable[x+p.hashtable[y]]
}

// Grad 勾配
func (p *PerlinNoise) Grad(hash int, a, b float64) float64 {
	switch hash % 0x04 {
	case 0x0:
		return a
	case 0x1:
		return -a
	case 0x2:
		return -b
	case 0x3:
		return b
	}

	return 0
}

func (p *PerlinNoise) PerlinNoise(x, y float64) float64 {
	xi := int(math.Floor(x))
	yi := int(math.Floor(y))
	xf := x - float64(xi)
	yf := y - float64(yi)

	a00 := p.Grad(p.getHash(xi, yi), xf, yf)
	a10 := p.Grad(p.getHash(xi+1, yi), xf-1, yf)
	a01 := p.Grad(p.getHash(xi, yi+1), xf, yf-1)
	a11 := p.Grad(p.getHash(xi+1, yi+1), xf-1, yf-1)

	xf = p.Fade(xf)
	yf = p.Fade(yf)

	return (p.Lerp(
		p.Lerp(a00, a10, xf),
		p.Lerp(a01, a11, xf),
		yf,
	) + 1) / 2
}

func (p *PerlinNoise) OctavePerlinNoise(x, y int) float64 {
	var (
		a          float64 = 1.
		f          float64 = 1.
		maxValue   float64 = 0.
		totalValue float64 = 0.
		per        float64 = 0.5
	)
	for i := 0; i < 5; i++ {
		totalValue += a * p.PerlinNoise(float64(x)*f, float64(y)*f)
		maxValue += a
		a *= per
		f *= 2
	}

	return math.Abs(totalValue / maxValue)
}

func setField() []color.Color {
	p := &PerlinNoise{}
	p.setHash(rand.Int63())
	field := make([]color.Color, WIDTH*HEIGHT)
	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			field[y*HEIGHT+x] = color.RGBA{
				uint8(255 * p.OctavePerlinNoise(x, y)),
				uint8(255 * p.OctavePerlinNoise(x, y)),
				uint8(255 * p.OctavePerlinNoise(x, y)),
				255,
			}
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
