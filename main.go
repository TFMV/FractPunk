package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math/cmplx"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Color palette
var colorPalette = []color.RGBA{
	{255, 69, 0, 255},   // Red
	{255, 215, 0, 255},  // Gold
	{138, 43, 226, 255}, // Blue
	{0, 255, 127, 255},  // Green
	{255, 20, 147, 255}, // Pink
}

func main() {
	rand.Seed(time.Now().UnixNano())
	const (
		width, height          = 1024, 1024
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
	)

	// Get GPT-4 generated text
	gptText, err := getGPT4Text()
	if err != nil {
		panic(err)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			img.Set(px, py, mandelbrot(z))
		}
	}

	// Add random annotations and flare
	addFlare(img, gptText)

	// Save the image to a file
	f, err := os.Create("fractpunk_fractal.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	var v complex128
	perturbation := complex(rand.Float64()*0.1-0.05, rand.Float64()*0.1-0.05)
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z + perturbation
		if cmplx.Abs(v) > 2 {
			return colorPalette[rand.Intn(len(colorPalette))]
		}
	}
	return color.Black
}

// addFlare adds random elements to the image
func addFlare(img *image.RGBA, text string) {
	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	randomAnnotations := []string{"Fnord!", "Kallisti!", "Ewige Blumenkraft", "Hail Eris!", "5 Tons of Flax", text}

	for i := 0; i < 5; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		col := colorPalette[rand.Intn(len(colorPalette))]
		addRandomShapes(img, x, y, col)
	}

	annotation := randomAnnotations[rand.Intn(len(randomAnnotations))]
	drawText(img, annotation, rand.Intn(width-100), rand.Intn(height-50), color.White)
}

// addRandomShapes adds random shapes to the image
func addRandomShapes(img *image.RGBA, x, y int, col color.Color) {
	for i := 0; i < 50; i++ {
		rx := x + rand.Intn(20) - 10
		ry := y + rand.Intn(20) - 10
		if rx >= 0 && ry >= 0 && rx < img.Bounds().Max.X && ry < img.Bounds().Max.Y {
			img.Set(rx, ry, col)
		}
	}
}

// drawText is a dummy function to illustrate random text placement
func drawText(img *image.RGBA, text string, x, y int, col color.Color) {
	fmt.Printf("Drawing text '%s' at (%d, %d)\n", text, x, y)
	// This function can be expanded with actual text drawing logic
	for i := 0; i < 50; i++ {
		img.Set(x+i, y, col)
	}
}

// getGPT4Text calls the GPT-4 API to get random text
func getGPT4Text() (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	apiKey := "YOUR_OPENAI_API_KEY"

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a whimsical phrase generator."},
			{"role": "user", "content": "Generate a random short whimsical phrase:"},
		},
		"max_tokens":  16,
		"temperature": 0.7,
	})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// Debugging: Print the API response
	fmt.Printf("API Response: %s\n", string(body))

	// Check if "choices" exists and is a non-empty array
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("unexpected API response format")
	}

	// Extract the text from the first choice
	text, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected API response format")
	}

	return text, nil
}
