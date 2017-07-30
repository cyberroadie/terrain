package main

import heightmap2 "terrain/heightmap"

func main() {
	heightmap := heightmap2.New(2048,2048)
	heightmap.DrawLine(50)
	heightmap.SaveImage("heightmap.png")
}
