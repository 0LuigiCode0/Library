package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/0LuigiCode0/library/image_filter"
)

func main() {
	path := "test5.png"
	newPath := "test5.min.png"
	levelDenoize := 2

	img, err := os.Open(path)
	if err != nil {
		fmt.Printf("image cannot open : %v\n", err)
		return
	}
	defer img.Close()

	oldMgg, f, err := image.Decode(img)
	if err != nil {
		fmt.Printf("image read is failed : %v\n", err)
		return
	}
	if oldMgg.Bounds().Dx() < 300 && oldMgg.Bounds().Dy() < 300 {
		fmt.Printf("image ia small\n")
		return
	}

	imgMin, err := os.OpenFile(newPath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Printf("min_image create is failed : %v\n", err)
		return
	}
	defer imgMin.Close()

	var height, width int
	coeffX := float32(oldMgg.Bounds().Dx()) / 400
	coeffY := float32(oldMgg.Bounds().Dy()) / 400
	ratio := coeffX / coeffY
	if coeffX > coeffY {
		coeffY *= ratio
		width = 400
		height = int(400 / ratio)
	} else {
		coeffX *= ratio
		width = int(400 / ratio)
		height = 400
	}
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newImage.Set(x, y, oldMgg.At(int(float32(x)*coeffX), int(float32(y)*coeffY)))
		}
	}

	image_filter.Denoize(newImage, levelDenoize)

	if f == "png" {
		if err = png.Encode(imgMin, newImage); err != nil {
			os.Remove(newPath)
			fmt.Printf("min_image create is failed : %v\n", err)
			return
		}
	} else {
		if err = jpeg.Encode(imgMin, newImage, &jpeg.Options{Quality: 100}); err != nil {
			os.Remove(newPath)
			fmt.Printf("min_image create is failed : %v\n", err)
			return
		}
	}
}
