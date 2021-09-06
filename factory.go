// Copyright (c) 2021 Gwenaël LE BODIC

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/fogleman/ln/ln"
)

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 0, 64)
}

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

	/*
		// vertical strips
		for i := 0; i <= c.Stripes; i++ {
			p := float64(i) / float64(c.Stripes)
			x := x1 + (x2-x1)*p
			y := y1 + (y2-y1)*p
			paths = append(paths, ln.Path{{x, y1, z1}, {x, y1, z2}})
			paths = append(paths, ln.Path{{x, y2, z1}, {x, y2, z2}}) // Vertical North
			paths = append(paths, ln.Path{{x1, y, z1}, {x1, y, z2}})
			paths = append(paths, ln.Path{{x2, y, z1}, {x2, y, z2}}) // Vertical North
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
	*/
	// vertical strips
	beam_size := 0.03
	//	println("stripes=" + fmt.Sprintf("%i", c.Stripes))
	for i := 0; i <= c.Stripes; i++ {
		p := float64(i) / float64(c.Stripes)
		x := x1 + (x2-x1)*p
		y := y1 + (y2-y1)*p

		if i == 0 || i == c.Stripes { // no beam / edges
			paths = append(paths, ln.Path{{x, y1, z1}, {x, y1, z2}})
			paths = append(paths, ln.Path{{x, y2, z1}, {x, y2, z2}}) // Vertical North 1
			paths = append(paths, ln.Path{{x1, y, z1}, {x1, y, z2}})
			paths = append(paths, ln.Path{{x2, y, z1}, {x2, y, z2}}) // Vertical North
		} else { // inner beam
			paths = append(paths, ln.Path{{x - beam_size, y1, z1}, {x - beam_size, y1, z2}})
			paths = append(paths, ln.Path{{x + beam_size, y1, z1}, {x + beam_size, y1, z2}})

			paths = append(paths, ln.Path{{x - beam_size, y2, z1}, {x - beam_size, y2, z2}}) // Vertical North 1
			paths = append(paths, ln.Path{{x + beam_size, y2, z1}, {x + beam_size, y2, z2}}) // Vertical North 2

			paths = append(paths, ln.Path{{x1, y - beam_size, z1}, {x1, y - beam_size, z2}})
			paths = append(paths, ln.Path{{x1, y + beam_size, z1}, {x1, y + beam_size, z2}})

			paths = append(paths, ln.Path{{x2, y - beam_size, z1}, {x2, y - beam_size, z2}}) // Vertical North
			paths = append(paths, ln.Path{{x2, y + beam_size, z1}, {x2, y + beam_size, z2}}) // Vertical North

		}

	}

	// horizontal strips

	var nbVstrips = int((z2 - z1) / 2 * 4)
	gap_length := float64(x2-x1) / float64(c.Stripes)
	for i := 0; i <= nbVstrips; i++ {
		p := float64(i) / float64(nbVstrips)
		z := z1 + (z2-z1)*p

		if i == 0 || i == nbVstrips {
			paths = append(paths, ln.Path{{x1, y1, z}, {x1, y2, z}})
			paths = append(paths, ln.Path{{x2, y1, z}, {x2, y2, z}})
			paths = append(paths, ln.Path{{x1, y1, z}, {x2, y1, z}})
			paths = append(paths, ln.Path{{x1, y2, z}, {x2, y2, z}})
		} else {

			paths = append(paths, ln.Path{{x1, y1, z}, {x1, y1 + gap_length - beam_size, z}})
			paths = append(paths, ln.Path{{x1, y2 - gap_length + beam_size, z}, {x1, y2, z}})
			for i := 1; i <= c.Stripes-1; i++ {
				p := float64(i) / float64(c.Stripes)
				y := y1 + (y2-y1)*p
				paths = append(paths, ln.Path{{x1, y + beam_size, z}, {x1, y + gap_length - beam_size, z}})
			}

			paths = append(paths, ln.Path{{x2, y1, z}, {x2, y1 + gap_length - beam_size, z}})
			paths = append(paths, ln.Path{{x2, y2 - gap_length + beam_size, z}, {x2, y2, z}})
			for i := 1; i <= c.Stripes-1; i++ {
				p := float64(i) / float64(c.Stripes)
				y := y1 + (y2-y1)*p
				paths = append(paths, ln.Path{{x2, y + beam_size, z}, {x2, y + gap_length - beam_size, z}})
			}

			paths = append(paths, ln.Path{{x1, y1, z}, {x1 + gap_length - beam_size, y1, z}})
			paths = append(paths, ln.Path{{x2 - gap_length + beam_size, y1, z}, {x2, y1, z}})
			for i := 1; i <= c.Stripes-1; i++ {
				p := float64(i) / float64(c.Stripes)
				x := x1 + (x2-x1)*p
				paths = append(paths, ln.Path{{x + beam_size, y1, z}, {x + gap_length - beam_size, y1, z}})
			}

			//paths = append(paths, ln.Path{{x1, y2, z}, {x2, y2, z}})
			paths = append(paths, ln.Path{{x1, y2, z}, {x1 + gap_length - beam_size, y2, z}})
			paths = append(paths, ln.Path{{x2 - gap_length + beam_size, y2, z}, {x2, y2, z}})
			for i := 1; i <= c.Stripes-1; i++ {
				p := float64(i) / float64(c.Stripes)
				x := x1 + (x2-x1)*p
				paths = append(paths, ln.Path{{x + beam_size, y2, z}, {x + gap_length - beam_size, y2, z}})
			}
		}
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
	eye_height := globalHeight + 2
	eye := ln.Vector{-20, -20, eye_height}
	center := ln.Vector{45, 45, 5}
	up := ln.Vector{0, 0, 1}

	// define rendering parameters
	width := 800.0   // rendered width
	height := 1600.0 // rendered height
	fovy := 60.0     // vertical field of view, degrees
	znear := 0.1     // near z plane
	zfar := 100.0    // far z plane
	step := 0.01     // how finely to chop the paths for visibility testing

	start := time.Now()

	paths := scene.Render(eye, center, up, width, height, fovy, znear, zfar, step)

	var filename = "out2/city" + strconv.Itoa(int(seed)) + "_" + FloatToString(eye.X) + "_" + FloatToString(eye.Y) + "_" + FloatToString(eye.Z) + "_" + FloatToString(center.X) + "_" + FloatToString(center.Y) + "_" + FloatToString(center.Z) + "_" + FloatToString(width) + "_" + FloatToString(height) + "_A"

	paths.WriteToPNG(filename+".png", width, height)
	paths.WriteToSVG(filename+".svg", width, height)
	elapsed := time.Since(start)

	println("A view generated.")
	println("Duration=" + elapsed.String())

	// B VIEW
	eye_height = 5.0 + globalHeight
	eye = ln.Vector{-5, -5, eye_height}
	center = ln.Vector{0, 0, 10}
	up = ln.Vector{0, 0, 1}

	// define rendering parameters
	width = 1024.0  // rendered width
	height = 1024.0 // rendered height
	fovy = 60.0     // vertical field of view, degrees
	znear = 0.1     // near z plane
	zfar = 100.0    // far z plane
	step = 0.01     // how finely to chop the paths for visibility testing

	start = time.Now()
	paths = scene.Render(eye, center, up, width, height, fovy, znear, zfar, step)

	//	paths := scene.Render(eye, center, up, width, height, 100, 0.1, 100, 0.01)
	filename = "out2/city" + strconv.Itoa(int(seed)) + "_" + FloatToString(eye.X) + "_" + FloatToString(eye.Y) + "_" + FloatToString(eye.Z) + "_" + FloatToString(center.X) + "_" + FloatToString(center.Y) + "_" + FloatToString(center.Z) + "_" + FloatToString(width) + "_" + FloatToString(height) + "_B"

	paths.WriteToPNG(filename+".png", width, height)
	paths.WriteToSVG(filename+".svg", width, height)

	elapsed = time.Since(start)
	println("B view generated.")
	println("Duration=" + elapsed.String())

}
