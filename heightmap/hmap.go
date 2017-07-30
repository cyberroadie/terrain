package heightmap

import (
	"image"
	"os"
	"image/png"
	"image/color"
	"image/draw"
	"log"
)

type HeightMap struct {
	img    *image.Gray16
}

const MAX_GREY  = ^uint16(0)

func New(height, width int) *HeightMap {
	img := image.NewGray16(image.Rect(0, 0, width, height))
	h := &HeightMap{img: img}
	h.setInitialHeight(0.5)
	return h
}

func (h *HeightMap) SaveImage(fileName string) {
	f, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, h.img)
}

func (h *HeightMap) setInitialHeight(height float64) {
	draw.Draw(h.img, h.img.Bounds(), &image.Uniform{getColor(height)}, image.ZP, draw.Src)
}

func (h *HeightMap) DrawLine(lineWidth int) {
	for i := 0; i < lineWidth; i++ {
		h.drawSingleLine(h.img.Bounds().Max.X / 2 - lineWidth /2 + i)
	}
}

func (h *HeightMap) drawSingleLine(offset int) {
	for y := h.img.Bounds().Min.Y; y < h.img.Bounds().Max.Y; y++ {
		h.img.SetGray16(offset, y,getColor(1.0))
	}
}

func getColor(c float64) color.Gray16 {
	if c < 0 || c > 1.0 {
		log.Fatalf("c value should be betweeen 0.0 and 1.0: %v", c)
	}
	intColor := uint16(float64(MAX_GREY) * c)
	return color.Gray16{intColor}
}