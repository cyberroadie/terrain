package heightmap

import (
	"image"
	"os"
	"image/png"
	"image/color"
	"image/draw"
	"log"
	"math"
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

func (h *HeightMap) DrawSingleLine(lineWidth int) {
	for i := 0; i < lineWidth; i++ {
		h.drawSingleLine(h.img.Bounds().Max.X / 2 - lineWidth /2 + i)
	}
}

func (h *HeightMap) drawSingleLine(offset int) {
	for y := h.img.Bounds().Min.Y; y < h.img.Bounds().Max.Y; y++ {
		h.img.SetGray16(offset, y,getColor(1.0))
	}
}

func (h *HeightMap) DrawLine(x1, y1, x2, y2 int, height float64, radius int) {
		h.Bresenham(x1, y1, x2, y2, radius, height)
}

func (h *HeightMap) Bresenham(x1, y1, x2, y2, radius int, height float64) {
	var dx, dy, e, slope int

	// Because drawing p1 -> p2 is equivalent to draw p2 -> p1,
	// I sort points in x-axis order to handle only half of possible cases.
	if x1 > x2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	dx, dy = x2-x1, y2-y1
	// Because point is x-axis ordered, dx cannot be negative
	if dy < 0 {
		dy = -dy
	}

	switch {

	// Is line a point ?
	case x1 == x2 && y1 == y2:
		h.drawCircle(x1, y1, radius, height)

		// Is line an horizontal ?
	case y1 == y2:
		for ; dx != 0; dx-- {
			h.drawCircle(x1, y1, radius, height)
			x1++
		}
		h.drawCircle(x1, y1, radius, height)

		// Is line a vertical ?
	case x1 == x2:
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for ; dy != 0; dy-- {
			h.drawCircle(x1, y1, radius, height)
			y1++
		}
		h.drawCircle(x1, y1, radius, height)

		// Is line a diagonal ?
	case dx == dy:
		if y1 < y2 {
			for ; dx != 0; dx-- {
				h.drawCircle(x1, y1, radius, height)
				x1++
				y1++
			}
		} else {
			for ; dx != 0; dx-- {
				h.drawCircle(x1, y1, radius, height)
				x1++
				y1--
			}
		}
		h.drawCircle(x1, y1, radius, height)

		// wider than high ?
	case dx > dy:
		if y1 < y2 {
			// BresenhamDxXRYD(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				h.drawCircle(x1, y1, radius, height)
				x1++
				e -= dy
				if e < 0 {
					y1++
					e += slope
				}
			}
		} else {
			// BresenhamDxXRYU(img, x1, y1, x2, y2, col)
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				h.drawCircle(x1, y1, radius, height)
				x1++
				e -= dy
				if e < 0 {
					y1--
					e += slope
				}
			}
		}
		h.drawCircle(x2, y2, radius, height)

		// higher than wide.
	default:
		if y1 < y2 {
			// BresenhamDyXRYD(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				h.drawCircle(x1, y1, radius, height)
				y1++
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		} else {
			// BresenhamDyXRYU(img, x1, y1, x2, y2, col)
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				h.drawCircle(x1, y1, radius, height)
				y1--
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		}
		h.drawCircle(x2, y2, radius, height)
	}
}

func (h *HeightMap) drawCircle(x1, y1, radius int, height float64) {
	for y :=-radius; y<=radius; y++ {
		for  x := -radius; x<=radius; x++ {
			if (x*x+y*y <= radius*radius) {
				h.img.SetGray16(x1 + x, y1 + y, getColor(height))
			}
		}
	}
}


func (h *HeightMap) drawLine(x1, y1, x2, y2 int, height float64) {
	var slope float64
	if x2 - x1 == 0 {
		slope = 0
	} else {
		slope = float64((y2 - y1)) / float64((x2 - x1))
	}
	b := float64(y1) - slope * float64(x1)

	for x := x1; x < x2; x++ {
		y := slope * float64(x) + b
		h.img.SetGray16(x, int(y), getColor(height))
	}

}

func getColor(c float64) color.Gray16 {
	if c < 0 || c > 1.0 {
		log.Fatalf("c value should be betweeen 0.0 and 1.0: %v", c)
	}
	intColor := uint16(float64(MAX_GREY) * c)
	return color.Gray16{intColor}
}

func (h *HeightMap) DrawHex(x1, y1, size, radiusLine int, height float64) {
	xs, ys := hexCorner(x1, y1, size, 1)

	for i := 2; i < 8  ; i++ {
		println("====")
		println(i)
		xe, ye := hexCorner(x1, y1, size, i)
		println(xs)
		println(ys)
		println(xe)
		println(ye)
		h.DrawLine(xs, ys, xe, ye, height, radiusLine)
		xs, ys = xe, ye

	}
}

func hexCorner(x1, y1, size, i int) (x, y int){
	angle_deg := float64(60 * i   + 30)
	angle_rad := math.Pi / 180 * angle_deg
	x = x1 + int(float64(size) * math.Cos(angle_rad))
	y = y1 +int(float64(size)* math.Sin(angle_rad))
	return
}

//function hex_corner(center, size, i):
//var angle_deg = 60 * i   + 30
//var angle_rad = PI / 180 * angle_deg
//return Point(center.x + size * cos(angle_rad),
//center.y + size * sin(angle_rad))
