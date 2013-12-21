package main

import (
	"fmt"
	gl "github.com/chsc/gogl/gl21"
	"github.com/jteeuwen/glfw"
	"os"
)

func RunGame(title string, width, height int, init, draw func()) {
	if err := glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}
	defer glfw.Terminate()

	glfw.OpenWindowHint(glfw.WindowNoResize, 1)

	if err := glfw.OpenWindow(width, height, 0, 0, 0, 0, 16, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "glfw: %s\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle(title)

	if err := gl.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "gl: %s\n", err)
	}

	if err := initScene(width, height, init); err != nil {
		fmt.Fprintf(os.Stderr, "init: %s\n", err)
		return
	}
	// defer destroyScene()

	for glfw.WindowParam(glfw.Opened) == 1 {
		drawScene(draw)
		glfw.SwapBuffers()
	}
}

func initScene(width, height int, init func()) (err error) {
	gl.Disable(gl.DEPTH_TEST)

	gl.ClearColor(0.5, 0.5, 0.5, 0.0)

	gl.Viewport(0, 0, gl.Sizei(width), gl.Sizei(height))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, gl.Double(width), gl.Double(height), 0, 0, 1)
	gl.MatrixMode(gl.MODELVIEW)

	init()

	return
}

/*
func destroyScene() {
  gl.DeleteTextures(1, &texture)
}
*/

func drawScene(draw func()) {
	// gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	draw()
	gl.Finish()
}
