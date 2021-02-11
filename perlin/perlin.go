package perlin

import (
	"math"
	"math/rand"
)

func New(width int, height int, seed int64, octerves int, persistence float64) *PerlinNoise {
	p := &PerlinNoise{
		octerves:    octerves,
		persistence: persistence,
		width:       width,
		height:      height,
	}

	p.setHash()
	return p
}

type PerlinNoise struct {
	hashtable   []int
	octerves    int
	persistence float64
	width       int
	height      int
	seed        int64
}

// Fade 6t^5 - 15t^4 + 10t^3.
func (p *PerlinNoise) Fade(t float64) float64 {
	return (6 * math.Pow(t, 5)) - (15 * math.Pow(5, 4)) + (10 * math.Pow(t, 3))
}

// Lerp 線形補間
func (p *PerlinNoise) Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func (p *PerlinNoise) setHash() {
	p.hashtable = make([]int, p.width*p.height)
	randtable := make([]int, p.height)
	for i := 0; i < p.height; i++ {
		randtable[i] = rand.Intn(p.height)
	}
	rand.Seed(p.seed)
	for i := 0; i < p.width*p.height; i++ {
		p.hashtable[i] = randtable[i%p.height]
	}
}

func (p *PerlinNoise) getHash(x, y int) int {
	x %= p.width
	y %= p.height
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
	xff := p.Fade(xf)
	yff := p.Fade(yf)

	xi %= p.width
	for xi < 0 {
		xi += p.width
	}
	yi %= p.height
	for yi < 0 {
		yi += p.height
	}

	a00 := p.Grad(p.getHash(xi, yi), xf, yf)
	a10 := p.Grad(p.getHash(xi+1, yi), xf-1, yf)
	a01 := p.Grad(p.getHash(xi, yi+1), xf, yf-1)
	a11 := p.Grad(p.getHash(xi+1, yi+1), xf-1, yf-1)

	return (p.Lerp(
		p.Lerp(a00, a10, xff),
		p.Lerp(a01, a11, xff),
		yff,
	) + 1) / 2
}

func (p *PerlinNoise) OctavePerlinNoise(x, y int) float64 {
	var (
		a          float64 = 1.
		f          float64 = 1.
		maxValue   float64 = 0.
		totalValue float64 = 0.
	)
	for i := 0; i < p.octerves; i++ {
		totalValue += a * p.PerlinNoise(float64(x)*f, float64(y)*f)
		maxValue += a
		a *= p.persistence
		f *= 2
	}

	return totalValue / maxValue
}
