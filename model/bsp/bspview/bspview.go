package main

import (
	"errors"
	"flag"
	gl "github.com/chsc/gogl/gl21"
	"github.com/ftrvxmtrx/groke/model/bsp"
	"github.com/jteeuwen/glfw"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"runtime/pprof"
)

const (
	Title  = "bspview"
	Width  = 1280
	Height = 800
)

const (
	moveW = 1 << iota
	moveA
	moveS
	moveD
	moveSpace
	moveCtrl
	turnLeft
	turnRight
	turnDown
	turnUp
)

var (
	move                int
	viewAngles, viewOrg [3]float64
	model               *bsp.Model
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if err := readModel(); err != nil {
		log.Fatal(err)
	}

	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	glfw.OpenWindowHint(glfw.WindowNoResize, 1)

	if err := glfw.OpenWindow(Width, Height, 0, 0, 0, 0, 16, 0, glfw.Windowed); err != nil {
		log.Fatal(err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle(Title)
	glfw.SetKeyCallback(onKey)

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	initScene()
	defer destroyScene()

	for glfw.WindowParam(glfw.Opened) == 1 {
		applyMove()
		drawScene()
		glfw.SwapBuffers()
	}
}

func initScene() {
	gl.Disable(gl.TEXTURE_2D)
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.ALPHA_TEST)
	gl.Enable(gl.LINE_SMOOTH)
	gl.Hint(gl.LINE_SMOOTH_HINT, gl.NICEST)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	gl.LineWidth(1.0)

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1)
	//gl.DepthFunc(gl.LEQUAL)

	gl.Viewport(0, 0, Width, Height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	perspective(110.0, 1.0, 4, 8192)
}

func perspective(fov, aspect, zNear, zFar gl.Double) {
	var xmin, xmax, ymin, ymax gl.Double

	ymax = zNear * gl.Double(math.Tan(float64(fov*math.Pi/360.0)))
	ymin = -ymax

	xmin = ymin * aspect
	xmax = ymax * aspect

	xmin += -1.0 / zNear
	xmax += -1.0 / zNear

	gl.Frustum(xmin, xmax, ymin, ymax, zNear, zFar)
}

func destroyScene() {
}

const moveSpeed = 5.0
const rotSpeed = 1.0

func cosGr(x float64) float64 {
	return moveSpeed * math.Cos(x*math.Pi/180.0)
}

func sinGr(x float64) float64 {
	return moveSpeed * math.Sin(x*math.Pi/180.0)
}

func applyMove() {
	if move&turnLeft != 0 {
		viewAngles[1] += rotSpeed
		if viewAngles[1] > 360.0 {
			viewAngles[1] -= 360.0
		}
	}

	if move&turnRight != 0 {
		viewAngles[1] -= rotSpeed
		if viewAngles[1] < 0.0 {
			viewAngles[1] = 360.0 - viewAngles[1]
		}
	}

	if move&turnDown != 0 {
		viewAngles[2] -= rotSpeed
		if viewAngles[2] < -90.0 {
			viewAngles[2] = -90.0
		}
	}

	if move&turnUp != 0 {
		viewAngles[2] += rotSpeed
		if viewAngles[2] > 90.0 {
			viewAngles[2] = 90.0
		}
	}

	if move&moveW != 0 {
		viewOrg[0] += sinGr(viewAngles[1])
		viewOrg[1] -= cosGr(viewAngles[1])
	}

	if move&moveA != 0 {
		viewOrg[0] += sinGr(viewAngles[1] + 90.0)
		viewOrg[1] -= cosGr(viewAngles[1] + 90.0)
	}

	if move&moveS != 0 {
		viewOrg[0] -= sinGr(viewAngles[1])
		viewOrg[1] += cosGr(viewAngles[1])
	}

	if move&moveD != 0 {
		viewOrg[0] += sinGr(viewAngles[1] - 90.0)
		viewOrg[1] -= cosGr(viewAngles[1] - 90.0)
	}

	if move&moveSpace != 0 {
		viewOrg[2] -= moveSpeed
	}

	if move&moveCtrl != 0 {
		viewOrg[2] += moveSpeed
	}
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Rotated(-90, 1, 0, 0)
	//gl.Rotated(90, 0, 0, 1)
	gl.Rotated(gl.Double(-viewAngles[2]), 1, 0, 0)
	gl.Rotated(gl.Double(-viewAngles[0]), 0, 1, 0)
	gl.Rotated(gl.Double(-viewAngles[1]), 0, 0, 1)
	gl.Translated(gl.Double(viewOrg[0]), gl.Double(viewOrg[1]), gl.Double(viewOrg[2]))

	var r, g, b gl.Ubyte

	for i, face := range model.Faces {
		r = gl.Ubyte(i) + 100
		g = gl.Ubyte(i>>1) + 100
		b = gl.Ubyte(i>>2) + 100

		if model.Triangle {
			gl.Begin(gl.TRIANGLES)
		} else {
			gl.Begin(gl.LINES)
		}

		gl.Color4ub(r, g, b, 0xff)
		for _, v := range face.Verts {
			gl.Vertex3d(gl.Double(v.Pos[0]), gl.Double(v.Pos[1]), gl.Double(v.Pos[2]))
		}

		gl.End()
	}
}

func readModel() error {
	if len(os.Args) < 2 {
		return errors.New("usage: bspview [options...] FILE")
	}

	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		return err
	}

	defer f.Close()
	model, err = bsp.Read(f, 0)
	if err == nil {
		log.Printf("numFaces: %d", len(model.Faces))
	}

	return err
}

func screenShot() {
	b := make([]uint8, Width*Height*4)
	gl.ReadPixels(0, 0, Width, Height, gl.RGBA, gl.UNSIGNED_BYTE, gl.Pointer(&b[0]))
	rect := image.Rect(0, 0, Width, Height)
	im := image.NRGBA{
		Pix:    b,
		Stride: Width * 4,
		Rect:   rect,
	}

	if f, err := os.Create("screenshot.png"); err == nil {
		defer f.Close()
		png.Encode(f, &im)
	}
}

func onKey(key, state int) {
	switch key {
	case glfw.KeyEsc:
		glfw.CloseWindow()
	case 'W':
		move ^= moveW
	case 'A':
		move ^= moveA
	case 'S':
		move ^= moveS
	case 'D':
		move ^= moveD
	case 'P':
		screenShot()
	case ' ':
		move ^= moveSpace
	case glfw.KeyLctrl:
		move ^= moveCtrl
	case glfw.KeyLeft:
		move ^= turnLeft
	case glfw.KeyRight:
		move ^= turnRight
	case glfw.KeyDown:
		move ^= turnDown
	case glfw.KeyUp:
		move ^= turnUp
	}
}
