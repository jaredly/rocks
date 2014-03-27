package main

import (
	gl "github.com/chsc/gogl/gl21"
	"github.com/jteeuwen/glfw"
)

const (
	Title  = "Rock Paper Scissors"
	Width  = 150
	Height = 150
	Scale  = 4
)

var matrix *[][]Point

func main() {
	RunGame(Title, Width*Scale, Height*Scale, func() {
		gl.PointSize(1.0 * Scale)

		x := 0
		y := 0
		dcolor := WHITE
		glfw.SetMouseButtonCallback(func(button, state int) {
			dcolor += 1
			if dcolor > BLUE {
				dcolor = WHITE
			}
		})
		glfw.SetMousePosCallback(func(mx, my int) {
			mx /= Scale
			my /= Scale
			if mx >= Width || my >= Height || mx < 0 || my < 0 {
				return
			}
			x = mx
			y = my
			if matrix != nil && dcolor != WHITE {
				(*matrix)[x][y].Color = dcolor
				(*matrix)[x][y].Intensity = 10
			}
		})
		glfw.SetKeyCallback(func(key, state int) {
			matrix = makeMatrix(Width, Height)
		})
		matrix = makeMatrix(Width, Height)
	}, func() {
		matrix = NextMatrix(Width, Height, matrix)
		gl.Begin(gl.POINTS)
		drawMatrix(matrix, Scale, Scale)
		gl.End()
	})
}
