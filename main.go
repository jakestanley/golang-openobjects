package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

type Room struct {
	cameraPosX float64
	cameraPosY float64
	cameraPosZ float64
	cameraRotY float64
}

type RoomObject struct {
	file      string
	locationX float64
	locationY float64
	locationZ float64
	scaleX    float64
	scaleY    float64
	scaleZ    float64
	rotationX float64
	rotationY float64
	rotationZ float64
}

var (
	assetsPath *string
)

func main() {

	assetsPath = flag.String("a", "./", "assets path")
	// flags must be parsed before they can be accessed
	flag.Parse()

	fmt.Println(assetsPath)

	a := app.App()
	// gui.Manager().Set(scene)

	// load an object
	// TODO clean path
	// assetSubPath := "Ceres.obj"
	// group := LoadModelToGroup(assetSubPath)
	// scene.Add(group)

	scenePath := fmt.Sprintf("%s/%s", *assetsPath, "cassandra_airlock.txt")
	scene, cam := ParseLoadScene(scenePath)

	// some code from the hellog3n demo. TODO: LICENSE

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}

func LoadModelToGroup(assetSubPath string) *core.Node {

	objPath := fmt.Sprintf("%s/%s", *assetsPath, assetSubPath)
	preprocessedPath := ObjPreprocessCache(objPath)

	decoded, err := obj.Decode(preprocessedPath, "")
	if err != nil {
		fmt.Println("error decoding obj: " + err.Error())
		os.Exit(1)
	}
	// TODO: this probably might end up using lots of memory...
	group, _ := decoded.NewGroup()

	return group
}

func ParseLoadScene(scenePath string) (*core.Node, *camera.Camera) {
	sceneFile, _ := os.Open(scenePath)
	scanner := bufio.NewScanner(sceneFile)

	scene := core.NewNode()

	// default camera for now. pull from room object
	// cam := camera.New(1)
	// cam.SetPosition(0, 0, 3)
	// scene.Add(cam)

	cam := camera.New(1)

	r_BeginRoom := regexp.MustCompile(`^begin room$`)
	r_BeginRoomObject := regexp.MustCompile(`^begin roomobject$`)
	for scanner.Scan() {
		line := scanner.Text()
		if r_BeginRoom.MatchString(line) {
			r := ParseRoom(scanner)
			cam.SetPosition(float32(r.cameraPosX), float32(r.cameraPosY), float32(r.cameraPosZ))
			cam.SetRotation(0, float32(r.cameraRotY), 0)
			scene.Add(cam)
		} else if r_BeginRoomObject.MatchString(line) {
			ro := ParseRoomObject(scanner)
			if len(ro.file) > 0 {
				group := LoadModelToGroup(ro.file)
				group.SetPosition(float32(ro.locationX), float32(ro.locationY), float32(ro.locationZ))
				group.SetScale(float32(ro.scaleX), float32(ro.scaleY), float32(ro.scaleZ))
				group.SetRotation(float32(ro.rotationX), float32(ro.rotationY), float32(ro.rotationZ))
				scene.Add(group)
			}
			fmt.Println(ro.file)
		}
	}
	scene.Add(helper.NewAxes(0.5))
	return scene, cam
}

func loadProps(scanner *bufio.Scanner, endOn *regexp.Regexp) map[string]string {
	props := map[string]string{}
	for scanner.Scan() {
		// todo: need to trim whitespace
		line := strings.TrimSpace(scanner.Text())
		if endOn.MatchString(line) {
			return props
		}
		kv := strings.Split(line, "=")
		if len(kv) == 2 {
			// ignore if we didn't split to two values. probably a better way to do it
			props[kv[0]] = kv[1]
			// fmt.Printf("%s: %s\n", kv[0], kv[1])
		}
	}
	// TODO throw error here as we have actually ended scanning without a terminator
	return props
}

func ParseRoom(scanner *bufio.Scanner) *Room {

	r := regexp.MustCompile(`^end room$`)
	props := loadProps(scanner, r)
	fmt.Println("Dumping room props")
	for k, v := range props {
		fmt.Printf("%s: %s\n", k, v)
	}

	cpx, _ := strconv.ParseFloat(props["cameraposx"], 32)
	cpy, _ := strconv.ParseFloat(props["cameraposy"], 32)
	cpz, _ := strconv.ParseFloat(props["cameraposz"], 32)
	cry, _ := strconv.ParseFloat(props["cameraroty"], 32)

	return &Room{
		cameraPosX: cpx,
		cameraPosY: cpy,
		cameraPosZ: cpz,
		cameraRotY: cry,
	}
}

func ParseRoomObject(scanner *bufio.Scanner) *RoomObject {

	r_End := regexp.MustCompile(`^end`)
	props := loadProps(scanner, r_End)
	fmt.Println("Dumping room object props")
	for k, v := range props {
		fmt.Printf("%s: %s\n", k, v)
	}

	lx, _ := strconv.ParseFloat(props["locationx"], 32)
	ly, _ := strconv.ParseFloat(props["locationy"], 32)
	lz, _ := strconv.ParseFloat(props["locationz"], 32)
	sx, _ := strconv.ParseFloat(props["scalex"], 32)
	sy, _ := strconv.ParseFloat(props["scaley"], 32)
	sz, _ := strconv.ParseFloat(props["scalez"], 32)
	rx, _ := strconv.ParseFloat(props["rotationx"], 32)
	ry, _ := strconv.ParseFloat(props["rotationy"], 32)
	rz, _ := strconv.ParseFloat(props["rotationz"], 32)

	return &RoomObject{
		file:      props["file"],
		locationX: lx,
		locationY: ly,
		locationZ: lz,
		scaleX:    sx,
		scaleY:    sy,
		scaleZ:    sz,
		rotationX: rx,
		rotationY: ry,
		rotationZ: rz,
	}
}

func ObjPreprocessCache(objPath string) string {

	processedPath := fmt.Sprintf("%s.pr", objPath)
	fmt.Printf("Looking for '%s'\n", processedPath)
	_, err := os.Stat(processedPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Processed file does not exist. Creating...")
			objFile, _ := os.Open(objPath)
			processedFile, _ := os.Create(processedPath)

			defer objFile.Close()
			defer processedFile.Close()

			r := regexp.MustCompile(`^s\s.`)

			scanner := bufio.NewScanner(objFile)
			for scanner.Scan() {
				line := scanner.Text()
				if !r.MatchString(line) {
					processedFile.WriteString(fmt.Sprintln(line))
				}
			}

		} else {
			panic(err)
		}
	} else {
		fmt.Println("Processed file exists. Loading...")
	}
	return processedPath
}
