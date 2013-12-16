package main

import (
	// "bytes"
	// "errors"
	"fmt"
	gl "github.com/chsc/gogl/gl21"
	"github.com/jteeuwen/glfw"
	// "image"
	// "image/png"
	// "io"
	"os"
)

const (
	Title  = "Spinning Gopher"
	Width  = 300
	Height = 300
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

	if err := glfw.OpenWindow(Width, Height, 0, 0, 0, 0, 16, 0, glfw.Windowed); err != nil {
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

	gl.Viewport(0, 0, Width, Height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, Width, Height, 0, 0, 1)
	gl.MatrixMode(gl.MODELVIEW)

	// gl.Enable(gl.POINT_SPRITE) // GL_POINT_SPRITE_ARB if you're
	// using the functionality as an extension.

	// gl.Enable(gl.POINT_SMOOTH)
	// gl.Enable(gl.BLEND)
	// gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.PointSize(1.0)

	/* assuming you have setup a 32-bit RGBA texture with a legal name */
	// gl.ActiveTexture(gl.TEXTURE0)
	// gl.Enable(gl.TEXTURE_2D)
	// gl.TexEnv(gl.POINT_SPRITE, gl.COORD_REPLACE, gl.TRUE);
	// gl.TexEnv(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.REPLACE);
	// gl.BindTexture(gl.TEXTURE_2D, texture)

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
	drawMatrix(matrix, 1, 1)
	gl.End()
	gl.Finish()
}
