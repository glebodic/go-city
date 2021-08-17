package main

import (
	"math/rand"
	"strconv"
	"time"

	//	"crypto/rand"

	"github.com/fogleman/ln/ln"
)

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
			scene.Add(&stripedcube)
			if rand.Float32() <= threshold_continuebuilding {
				newHeight := float64(2.0 + rand.Intn(10))
				buildExtension(scene, ln.Vector{(max.X-min.X)/2 + min.X, min.Y, max.Z}, ln.Vector{max.X, (max.Y-min.Y)/2 + min.Y, max.Z + newHeight})
			}
		}
	}
}

func main() {
	scene := ln.Scene{}
	area_min_x := -50.0
	area_min_y := -50.0
	area_max_x := 50.0
	area_max_y := 50.0
	level_height := 6.0

	var seed = time.Now().UnixNano()
	rand.Seed(seed)
	println("Seed = " + strconv.Itoa(int(seed)))

	// This is the building base
	strippedcube := createStripedCube(ln.Vector{area_min_x, area_min_y, 0}, ln.Vector{area_max_x, area_max_y, level_height}, int((area_max_x-area_min_x)*3))
	scene.Add(&strippedcube)

	// Now let's build next levels, recursively
	buildExtension(&scene, ln.Vector{area_min_x, area_min_y, level_height}, ln.Vector{area_max_x, area_max_y, 2 * level_height})

	//	View from bottom
	/*	eye := ln.Vector{-40, -40, 0}
		center := ln.Vector{0, 0, 30}
		up := ln.Vector{0, 0, 1} */

	// View from top
	eye := ln.Vector{-10, -10, 60}
	center := ln.Vector{20, 20, 20}
	up := ln.Vector{0, 0, 1}

	// define rendering parameters
	width := 4096.0  // rendered width
	height := 4096.0 // rendered height
	fovy := 60.0     // vertical field of view, degrees
	znear := 0.1     // near z plane
	zfar := 100.0    // far z plane
	step := 0.01     // how finely to chop the paths for visibility testing

	paths := scene.Render(eye, center, up, width, height, fovy, znear, zfar, step)

	//	paths := scene.Render(eye, center, up, width, height, 100, 0.1, 100, 0.01)
	var filename = "out/city" + strconv.Itoa(int(seed))

	paths.WriteToPNG(filename+".png", width, height)
	paths.WriteToSVG(filename+".svg", width, height)
	println("Done!")
	// paths.Print()
}