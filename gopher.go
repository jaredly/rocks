package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl21"
	"github.com/jteeuwen/glfw"
	"os"
)

const (
  Title  = "Spinning Gopher"
  Width  = 300
  Height = 300
  Scale = 1
)

var (
	matrix     *[][]Point
	texture    gl.Uint
	rotx, roty gl.Float
	ambient    []gl.Float = []gl.Float{0.5, 0.5, 0.5, 1}
	diffuse    []gl.Float = []gl.Float{1, 1, 1, 1}
	lightpos   []gl.Float = []gl.Float{-5, 5, 10, 0}
)

func main() {
	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}
	defer glfw.Terminate()

	glfw.OpenWindowHint(glfw.WindowNoResize, 1)

	if err := glfw.OpenWindow(Width*Scale, Height*Scale, 0, 0, 0, 0, 16, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle(Title)

	if err := gl.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "gl: %s\n", err)
	}

	if err := initScene(); err != nil {
		fmt.Fprintf(os.Stderr, "init: %s\n", err)
		return
	}
	defer destroyScene()

	for glfw.WindowParam(glfw.Opened) == 1 {
		StepMatrix(Width, Height, matrix)
		drawScene()
		glfw.SwapBuffers()
	}
}

func initScene() (err error) {
	gl.Disable(gl.DEPTH_TEST)

	gl.ClearColor(0.5, 0.5, 0.5, 0.0)

	gl.Viewport(0, 0, Width*Scale, Height*Scale)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, Width*Scale, Height*Scale, 0, 0, 1)
	gl.MatrixMode(gl.MODELVIEW)

	gl.PointSize(1.0*Scale)

	x := 0
	y := 0
	dcolor := WHITE
	glfw.SetMouseButtonCallback(func (button, state int) {
		dcolor += 1
		if dcolor > BLUE {
			dcolor = WHITE
		}
	})
	glfw.SetMousePosCallback(func (mx, my int) {
                mx /= Scale
                my /= Scale
                if mx >= Width || my >= Height || mx < 0 || my < 0 {
                    return
                }
		x = mx
		y = my
		if matrix != nil {
			(*matrix)[x][y].Color = dcolor
			(*matrix)[x][y].Intensity = 10
		}
	})

	matrix = makeMatrix(Width, Height)

	return
}

func destroyScene() {
	gl.DeleteTextures(1, &texture)
}

func drawScene() {
	// gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.Begin(gl.POINTS)
	drawMatrix(matrix, Scale, Scale)
	gl.End()
	gl.Finish()
}
