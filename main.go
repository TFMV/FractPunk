// FractPunk is a CLI tool that generates fractal images with artistic overlays.
// It chooses between Mandelbrot and Julia sets, and adds various artistic elements to the generated images.

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"math/rand"
	"os"
	"time"
)

// A palette used for drawing overlay shapes.
var colorPalette = []color.RGBA{
	{255, 69, 0, 255},   // Red
	{255, 215, 0, 255},  // Gold
	{138, 43, 226, 255}, // Blue
	{0, 255, 127, 255},  // Green
	{255, 20, 147, 255}, // Pink
}

func main() {
	// Replace rand.Seed with a local random source
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const (
		width, height          = 1024, 1024
		xmin, ymin, xmax, ymax = -2, -2, 2, 2
	)

	// Update rand.Intn calls to use rng.Intn
	useJulia := rng.Intn(2) == 0
	var c complex128
	if useJulia {
		// Choose a random constant for the Julia set.
		c = complex(rng.Float64()*2-1, rng.Float64()*2-1)
		fmt.Printf("Using Julia set with constant: %v\n", c)
	} else {
		fmt.Println("Using Mandelbrot set")
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			col := fractal(z, useJulia, c)
			img.Set(px, py, col)
		}
	}

	// Add artistic overlays: circles, sparkles, and a quirky annotation.
	addArtisticOverlays(img, "FractPunk")

	// Save the image using a random filename.
	fileName := fmt.Sprintf("%s.png", randomString(10))
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
	fmt.Printf("Saved artwork as %s\n", fileName)
}

// fractal computes a fractal color based on point z.
// When useJulia is true it renders a Julia set (using constant c),
// otherwise it computes the Mandelbrot set.
func fractal(z complex128, useJulia bool, c complex128) color.Color {
	const iterations = 300
	var v complex128
	if useJulia {
		v = z
	} else {
		v = 0
	}
	var n uint16
	for n = 0; n < iterations; n++ {
		if useJulia {
			v = v*v + c
		} else {
			v = v*v + z
		}
		if cmplx.Abs(v) > 2 {
			break
		}
	}
	if n == iterations {
		return color.Black
	}
	// Apply smoothing for a continuous gradient.
	smooth := float64(n) - math.Log(math.Log(cmplx.Abs(v)))/math.Log(2)
	t := smooth / iterations
	// Map the normalized value to a hue (0-360°), with high saturation.
	return hsvToRGB(360*t, 0.8, 1)
}

// hsvToRGB converts an HSV color (h in degrees, s and v in [0,1]) to an RGBA color.
func hsvToRGB(h, s, v float64) color.RGBA {
	c := v * s
	hPrime := h / 60.0
	x := c * (1 - math.Abs(math.Mod(hPrime, 2)-1))
	m := v - c
	var r, g, b float64
	switch {
	case hPrime < 1:
		r, g, b = c, x, 0
	case hPrime < 2:
		r, g, b = x, c, 0
	case hPrime < 3:
		r, g, b = 0, c, x
	case hPrime < 4:
		r, g, b = 0, x, c
	case hPrime < 5:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	return color.RGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
}

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = 'a' + byte(rand.Intn(26))
	}
	return string(b)
}

// drawGradientCircle renders a circle with a smooth color transition.
func drawGradientCircle(img *image.RGBA, cx, cy, r int, startCol, endCol color.RGBA) {
	for i := r; i > 0; i-- {
		t := float64(i) / float64(r)
		col := interpolateColor(startCol, endCol, t)
		drawCircle(img, cx, cy, i, col)
	}
}

// interpolateColor smoothly transitions between two colors.
func interpolateColor(c1, c2 color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(c1.R)*(1-t) + float64(c2.R)*t),
		G: uint8(float64(c1.G)*(1-t) + float64(c2.G)*t),
		B: uint8(float64(c1.B)*(1-t) + float64(c2.B)*t),
		A: 255,
	}
}

// addArtisticOverlays adds artistic elements to the fractal image.
func addArtisticOverlays(img *image.RGBA, baseText string) {
	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	randomAnnotations := []string{"Fnord!", "Kallisti!", "Ewige Blumenkraft", "Hail Eris!", "5 Tons of Flax", baseText}

	// Draw gradient circles for a smoother artistic look.
	for i := 0; i < 8; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		radius := rand.Intn(50) + 10
		startCol := colorPalette[rand.Intn(len(colorPalette))]
		endCol := colorPalette[rand.Intn(len(colorPalette))]
		drawGradientCircle(img, x, y, radius, startCol, endCol)
	}

	// Scatter sparkles with variable opacity.
	for i := 0; i < 60; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		alpha := uint8(rand.Intn(128) + 127) // Random transparency (127-255)
		col := colorPalette[rand.Intn(len(colorPalette))]
		img.Set(x, y, color.RGBA{col.R, col.G, col.B, alpha})
	}

	// Add chaotic artistic lines.
	for i := 0; i < 5; i++ {
		x1, y1 := rand.Intn(width), rand.Intn(height)
		x2, y2 := rand.Intn(width), rand.Intn(height)
		col := colorPalette[rand.Intn(len(colorPalette))]
		drawLine(img, x1, y1, x2, y2, col)
	}

	// Place an annotation dynamically.
	annotation := randomAnnotations[rand.Intn(len(randomAnnotations))]
	drawText(img, annotation, width/4+rand.Intn(width/2), height/4+rand.Intn(height/2), color.White)
}

// drawLine uses Bresenham's algorithm to draw a line.
func drawLine(img *image.RGBA, x0, y0, x1, y1 int, col color.Color) {
	dx := int(math.Abs(float64(x1 - x0)))
	dy := int(math.Abs(float64(y1 - y0)))
	sx, sy := -1, -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy

	for {
		if x0 >= 0 && x0 < img.Bounds().Max.X && y0 >= 0 && y0 < img.Bounds().Max.Y {
			img.Set(x0, y0, col)
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// drawCircle renders a circle centered at (cx, cy) with radius r.
func drawCircle(img *image.RGBA, cx, cy, r int, col color.Color) {
	steps := 360
	for i := 0; i < steps; i++ {
		theta := 2 * math.Pi * float64(i) / float64(steps)
		x := cx + int(float64(r)*math.Cos(theta))
		y := cy + int(float64(r)*math.Sin(theta))
		if x >= 0 && y >= 0 && x < img.Bounds().Max.X && y < img.Bounds().Max.Y {
			img.Set(x, y, col)
		}
	}
}

// drawText is a placeholder to indicate where and how text might be rendered.
func drawText(img *image.RGBA, text string, x, y int, col color.Color) {
	fmt.Printf("Drawing text '%s' at (%d, %d)\n", text, x, y)
	// For now, simply draw a horizontal line as a stand-in for text.
	for i := 0; i < len(text)*8; i++ {
		if x+i < img.Bounds().Max.X && y < img.Bounds().Max.Y {
			img.Set(x+i, y, col)
		}
	}
}
