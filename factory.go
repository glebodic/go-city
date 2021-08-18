package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/fogleman/ln/ln"
)

// Global variable
var globalHeight float64

type Antenna struct {
	ln.Cube
	V0 ln.Vector
}

func (a *Antenna) Paths() ln.Paths {
	var paths ln.Paths
	paths = append(paths, ln.Path{{0, 0, 0}, {a.V0.X, a.V0.Y, a.V0.Z}})

	return paths
}

func createAntenna(top ln.Vector) (c Antenna) {

	cube := ln.Cube{top, top, ln.Box{top, top}}
	antenna := Antenna{cube, top}
	return antenna
}

type StripedCube struct {
	ln.Cube
	Stripes int
}

func (c *StripedCube) Paths() ln.Paths {
	var paths ln.Paths
	x1, y1, z1 := c.Min.X, c.Min.Y, c.Min.Z
	x2, y2, z2 := c.Max.X, c.Max.Y, c.Max.Z

	// vertical strips
	for i := 0; i <= c.Stripes; i++ {
		p := float64(i) / float64(c.Stripes)
		x := x1 + (x2-x1)*p
		y := y1 + (y2-y1)*p
		paths = append(paths, ln.Path{{x, y1, z1}, {x, y1, z2}})
		paths = append(paths, ln.Path{{x, y2, z1}, {x, y2, z2}})
		paths = append(paths, ln.Path{{x1, y, z1}, {x1, y, z2}})
		paths = append(paths, ln.Path{{x2, y, z1}, {x2, y, z2}})
	}

	// horizontal strips

	var nbVstrips = int((z2 - z1) / 2 * 4)
	for i := 0; i <= nbVstrips; i++ {
		p := float64(i) / float64(nbVstrips)
		z := z1 + (z2-z1)*p
		paths = append(paths, ln.Path{{x1, y1, z}, {x1, y2, z}})
		paths = append(paths, ln.Path{{x2, y1, z}, {x2, y2, z}})
		paths = append(paths, ln.Path{{x1, y1, z}, {x2, y1, z}})
		paths = append(paths, ln.Path{{x1, y2, z}, {x2, y2, z}})
	}

	return paths
}

func createStripedCube(min ln.Vector, max ln.Vector, nbStripes int) (c StripedCube) {
	cube := ln.Cube{min, max, ln.Box{min, max}}
	stripedcube := StripedCube{cube, nbStripes}
	return stripedcube
}

func buildExtension(scene *ln.Scene, min ln.Vector, max ln.Vector) {

	if max.Z < 600 && (max.X-min.X > 1) {
		//		choice := rand.Float64()
		threshold_slicebuild := float32(0.5)
		threshold_continuebuilding := float32(0.8)

		newNbStripes := int((max.X - min.X) / 2 * 4)

		if rand.Float32() <= threshold_slicebuild {
			buildExtension(scene, ln.Vector{min.X, min.Y, min.Z}, ln.Vector{(max.X-min.X)/2 + min.X, (max.Y-min.Y)/2 + min.Y, max.Z})
		} else {
			stripedcube := createStripedCube(ln.Vector{min.X, min.Y, min.Z}, ln.Vector{(max.X-min.X)/2 + min.X, (max.Y-min.Y)/2 + min.Y, max.Z}, newNbStripes)
			if max.Z > globalHeight {
				globalHeight = max.Z
			}
			scene.Add(&stripedcube)
			if rand.Float32() <= threshold_continuebuilding {
				newHeight := float64(2.0 + rand.Intn(10))
				buildExtension(scene, ln.Vector{min.X, min.Y, max.Z}, ln.Vector{(max.X-min.X)/2 + min.X, (max.Y-min.Y)/2 + min.Y, max.Z + newHeight})
			}
		}

		if rand.Float32() <= threshold_slicebuild {
			buildExtension(scene, ln.Vector{(max.X-min.X)/2 + min.X, (max.Y-min.Y)/2 + min.Y, min.Z}, ln.Vector{max.X, max.Y, max.Z})
		} else {
			stripedcube := createStripedCube(ln.Vector{(max.X-min.X)/2 + min.X, (max.Y-min.Y)/2 + min.Y, min.Z}, ln.Vector{max.X, max.Y, max.Z}, newNbStripes)
			if max.Z > globalHeight {
				globalHeight = max.Z
			}
			scene.Add(&stripedcube)
			if rand.Float32() <= threshold_continuebuilding {
				newHeight := float64(2.0 + rand.Intn(10))
				buildExtension(scene, ln.Vector{(max.X-min.X)/2 + min.X, (max.Y-min.Y)/2 + min.Y, max.Z}, ln.Vector{max.X, max.Y, max.Z + newHeight})
			}
		}

		if rand.Float32() <= threshold_slicebuild {
			buildExtension(scene, ln.Vector{min.X, (max.Y-min.Y)/2 + min.Y, min.Z}, ln.Vector{(max.X-min.X)/2 + min.X, max.Y, max.Z})
		} else {
			stripedcube := createStripedCube(ln.Vector{min.X, (max.Y-min.Y)/2 + min.Y, min.Z}, ln.Vector{(max.X-min.X)/2 + min.X, max.Y, max.Z}, newNbStripes)
			if max.Z > globalHeight {
				globalHeight = max.Z
			}
			scene.Add(&stripedcube)
			if rand.Float32() <= threshold_continuebuilding {
				newHeight := float64(2.0 + rand.Intn(10))
				buildExtension(scene, ln.Vector{min.X, (max.Y-min.Y)/2 + min.Y, max.Z}, ln.Vector{(max.X-min.X)/2 + min.X, max.Y, max.Z + newHeight})
			}
		}

		if rand.Float32() <= threshold_slicebuild {
			buildExtension(scene, ln.Vector{(max.X-min.X)/2 + min.X, min.Y, min.Z}, ln.Vector{max.X, (max.Y-min.Y)/2 + min.Y, max.Z})
		} else {
			stripedcube := createStripedCube(ln.Vector{(max.X-min.X)/2 + min.X, min.Y, min.Z}, ln.Vector{max.X, (max.Y-min.Y)/2 + min.Y, max.Z}, newNbStripes)
			if max.Z > globalHeight {
				globalHeight = max.Z
			}
			scene.Add(&stripedcube)
			if rand.Float32() <= threshold_continuebuilding {
				newHeight := float64(2.0 + rand.Intn(10))
				buildExtension(scene, ln.Vector{(max.X-min.X)/2 + min.X, min.Y, max.Z}, ln.Vector{max.X, (max.Y-min.Y)/2 + min.Y, max.Z + newHeight})
			}
		}
	} else {
		if max.Z > 600 {
			println("Reached max height")
		}
	}
}

func main() {
	scene := ln.Scene{}
	area_min_x := -50.0
	area_min_y := -50.0
	area_max_x := 50.0
	area_max_y := 50.0
	level_height := 1.0
	globalHeight = 0

	var seed int64
	if len(os.Args[1:]) == 0 {
		seed = time.Now().UnixNano()
		println("New seed = " + strconv.Itoa(int(seed)))
	} else {
		seedArg := os.Args[1]
		i, err := strconv.Atoi(seedArg)
		seed = int64(i)
		println("Reused seed = " + strconv.Itoa(int(seed)))
		if err != nil {
			fmt.Println(err)
		}
	}
	rand.Seed(int64(seed))

	// This is the building base
	strippedcube := createStripedCube(ln.Vector{area_min_x, area_min_y, 0}, ln.Vector{area_max_x, area_max_y, level_height}, int((area_max_x-area_min_x)*3))
	scene.Add(&strippedcube)

	// Now let's build next levels, recursively
	buildExtension(&scene, ln.Vector{area_min_x, area_min_y, level_height}, ln.Vector{area_max_x, area_max_y, 2 * level_height})
	//oneLine := createAntenna(ln.Vector{50.0, 50.0, 50.0})
	//scene.Add(&oneLine)

	println("Global model height = " + fmt.Sprintf("%f", globalHeight))

	//	View from bottom
	/*	eye := ln.Vector{-40, -40, 0}
		center := ln.Vector{0, 0, 30}
		up := ln.Vector{0, 0, 1} */

	// A VIEW
	eye_height := 5.0 + globalHeight
	eye := ln.Vector{-20, -20, eye_height - 10.0}
	center := ln.Vector{25, 25, 20}
	up := ln.Vector{0, 0, 1}

	// define rendering parameters
	width := 4096.0  // rendered width
	height := 4096.0 // rendered height
	fovy := 60.0     // vertical field of view, degrees
	znear := 0.1     // near z plane
	zfar := 100.0    // far z plane
	step := 0.01     // how finely to chop the paths for visibility testing

	paths := scene.Render(eye, center, up, width, height, fovy, znear, zfar, step)

	var filename = "out/city" + strconv.Itoa(int(seed)) + "_A"

	paths.WriteToPNG(filename+".png", width, height)
	//	paths.WriteToSVG(filename+".svg", width, height)
	println("A view generated.")

	// B VIEW
	eye_height = 5.0 + globalHeight
	println("B global height=", eye_height)
	eye = ln.Vector{-5, -5, eye_height}
	center = ln.Vector{0, 0, 10}
	up = ln.Vector{0, 0, 1}

	// define rendering parameters
	width = 4096.0  // rendered width
	height = 4096.0 // rendered height
	fovy = 60.0     // vertical field of view, degrees
	znear = 0.1     // near z plane
	zfar = 100.0    // far z plane
	step = 0.01     // how finely to chop the paths for visibility testing

	paths = scene.Render(eye, center, up, width, height, fovy, znear, zfar, step)

	//	paths := scene.Render(eye, center, up, width, height, 100, 0.1, 100, 0.01)
	filename = "out/city" + strconv.Itoa(int(seed)) + "_B"

	paths.WriteToPNG(filename+".png", width, height)
	//paths.WriteToSVG(filename+".svg", width, height)
	println("B view generated.")

	// paths.Print()
}
