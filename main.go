package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
)

const (
	GTitle  = "Template"
	WWidth  = 800
	WHeight = 600
)

func main() {

	var assetsPath = flag.String("a", "./", "assets path")
	// flags must be parsed before they can be accessed
	flag.Parse()

	fmt.Println(assetsPath)

	a := app.App()
	scene := core.NewNode()
	// gui.Manager().Set(scene)

	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)
	camera.NewOrbitControl(cam)

	// load an object
	// TODO clean path
	objPath := fmt.Sprintf("%s/%s", *assetsPath, "Ceres.obj")
	decoded, err := obj.Decode(objPath, "")
	if err != nil {
		fmt.Println("error decoding obj: " + err.Error())
		os.Exit(1)
	}
	group, _ := decoded.NewGroup()
	group.SetScale(0.3, 0.3, 0.3)
	group.SetPosition(0, 0, 0)
	scene.Add(group)

	// some code from the hellog3n demo. TODO LICENSE

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
