package main

import heightmap2 "terrain/heightmap"

func main() {
	heightmap := heightmap2.New(2048,2048)
	//heightmap.DrawSingleLine(50)
	//heightmap.DrawLine(1024, 1536, 581, 1279, 1.0, 1)
	//heightmap.DrawLine(1024, 1024, 2048, 2048, 1.0, 10)
	//heightmap.DrawLine(2048, 2048, 1024, 1024, 0.7, 10)
	heightmap.DrawHex(1024, 1024, 512,10, 0.8)
	heightmap.SaveImage("heightmap.png")
}
