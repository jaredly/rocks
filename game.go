package main

import (
	gl "github.com/chsc/gogl/gl21"
  "math/rand"
  // "fmt"
)

type Point struct {
	Color int
	Intensity int
}

const (
  WHITE = 1
  RED = 2
  GREEN = 3
  BLUE = 4
  ORANGE = 5
)

type Color struct {
  R gl.Float
  G gl.Float
  B gl.Float
  A gl.Float
}

func UseColor(c Color) {
  gl.Color4f(c.R, c.G, c.B, c.A)
}

func PointColor(p Point) Color {
  n := gl.Float(float64(p.Intensity) / 10.0)
  /*
  if p.Color == WHITE {
    return Color{1.0, 1.0, 1.0, 1}
  }
  if p.Color == RED {
    return Color{1, 1-n, 1-n, 1}
  }
  if p.Color == GREEN {
    return Color{1-n, 1, 1-n, 1}
  }
  return Color{1-n, 1-n, 1, 1}
  */
  if p.Color == WHITE {
    return Color{0, 0, 0, 1}
  }
  if p.Color == RED {
    return Color{n, 0, 0, 1}
  }
  if p.Color == GREEN {
    return Color{0, n, 0, 1}
  }
  return Color{0, 0, n, 1}
}

func Dir() (int, int) {
  n := rand.Intn(8)
  d := [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}}[n]
  return d[0], d[1]
}

func Win(p *Point) {
  if p.Intensity >= 10 {
    return
  }
  p.Intensity += 1
}

func Lose(p *Point, ncolor int) {
  if p.Intensity > 1 {
    p.Intensity -= 1
    return
  }
  p.Intensity = 0
  p.Color = WHITE
}

func Beats(c1, c2 int) bool {
  if c1 == RED {
    return c2 == GREEN
  }
  if c1 == BLUE {
    return c2 == RED
  }
  return c2 == BLUE
}

func StepOne(x, y, width, height int, matrix *[][]Point) {
  dx, dy := Dir()
  x2 := x + dx
  y2 := y + dy
  if x2 < 0 || x2 >= width {
    return
  }
  if y2 < 0 || y2 >= height {
    return
  }
  p1 := &(*matrix)[x][y]
  p2 := &(*matrix)[x2][y2]
  if p1.Color == p2.Color {
    Lose(p1, WHITE)
    return
  }
  if p1.Color == WHITE {
    if p2.Intensity < 2 {
      return
    }
    p1.Intensity = p2.Intensity - 1
    p1.Color = p2.Color
    return
  }
  if p2.Color == WHITE {
    if p1.Intensity < 2 {
      return
    }
    p2.Intensity = p1.Intensity - 1
    p2.Color = p1.Color
    return
  }
  if Beats(p2.Color, p1.Color) {
    Win(p2)
    Lose(p1, p2.Color)
  } else {
    Win(p1)
    Lose(p2, p2.Color)
  }
}

func StepMatrix(width, height int, matrix *[][]Point) {
  for x := range *matrix {
    for y := range (*matrix)[x] {
      StepOne(x, y, width, height, matrix)
    }
  }
}

func makeMatrix(width, height int) *[][]Point {
  m := make([][]Point, width)
  for x := range m {
    m[x] = make([]Point, height)
    for y := range m[x] {
      m[x][y].Color = RED
      m[x][y].Intensity = 10
    }
  }

  for x := range m {
    z := (x/10) % 3
    if z == 1 {
      m[x][x].Color = RED
      m[x][height-x-1].Color = RED
    } else if z == 2 {
      m[x][x].Color = GREEN
      m[x][height-x-1].Color = GREEN
    } else {
      m[x][x].Color = BLUE
      m[x][height-x-1].Color = BLUE
    }
    m[x][x].Intensity = 10
    m[x][height-x-1].Intensity = 10
  }

  /*
  a := width/3
  drawStripe(a, a, 20, 20, &m, RED)
  drawStripe(a+20, a, 20, 20, &m, BLUE)
  drawStripe(a, a+20, 20, 20, &m, GREEN)
  drawStripe(a+20, a+20, 20, 20, &m, RED)
  a += 40
  drawStripe(a, a, 20, 20, &m, RED)
  drawStripe(a+20, a, 20, 20, &m, BLUE)
  drawStripe(a, a+20, 20, 20, &m, GREEN)
  drawStripe(a+20, a+20, 20, 20, &m, RED)

  drawStripe(a, a-20, 20, 20, &m, GREEN)
  drawStripe(a-20, a, 20, 20, &m, BLUE)
  */

  return &m
}

func drawStripe(x, y, w, h int, m *[][]Point, c int) {
  for a := x; a < x + w; a++ {
    for b := y; b < y + h; b++ {
      (*m)[a][b].Color = c
      (*m)[a][b].Intensity = 10
    }
  }
}

func drawMatrix(m *[][]Point, xby, yby gl.Float) {
  for x, row := range *m {
    for y, p := range row {
      UseColor(PointColor(p))
      // fmt.Println(p, PointColor(p))
      // gl.Color4f(1, 0, 0, 1)
      gl.Vertex2f(gl.Float(x) * xby, gl.Float(y) * yby)
    }
  }
}

