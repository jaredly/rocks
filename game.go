package main

import (
	gl "github.com/chsc/gogl/gl21"
	"math/rand"
	// "fmt"
)

type Point struct {
	Color     int
	Intensity int
}

const (
	WHITE  = 1
	RED    = 2
	GREEN  = 3
	BLUE   = 4
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

func Dirs() [][2]int {
	return [][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
}
func Diags() [][2]int {
	return [][2]int{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}}
}

func Dir() (int, int) {
	n := rand.Intn(8)
	d := Dirs()[n]
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
	p.Intensity = 1
	p.Color = ncolor
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

func safe(x, y, dx, dy, width, height int) bool {
	x2 := x + dx
	y2 := y + dy
	return x2 >=0 && y2 >= 0 && x2 < width && y2 < height
}

func StepOne(x, y, width, height int, p *Point, matrix *[][]Point) {
	r := 0.0
	rm := 0.0
	g := 0.0
	gm := 0.0
	b := 0.0
	bm := 0.0
	dirs := Dirs()
	for _, d := range dirs {
		if !safe(x,y,d[0],d[1],width,height) {
			continue
		}
		p2 := &(*matrix)[x+d[0]][y+d[1]]
		if p2.Color == RED {
			r += 1
			rm += float64(p2.Intensity)
		}
		if p2.Color == GREEN {
			g += 1
			gm += float64(p2.Intensity)
		}
		if p2.Color == BLUE {
			b += 1
			bm += float64(p2.Intensity)
		}
	}
	diags := Diags()
	for _, d := range diags {
		if !safe(x,y,d[0],d[1],width,height) {
			continue
		}
		p2 := &(*matrix)[x+d[0]][y+d[1]]
		if p2.Color == RED {
			r += .5
			rm += float64(p2.Intensity)/2.0
		}
		if p2.Color == GREEN {
			g += .5
			gm += float64(p2.Intensity)/2.0
		}
		if p2.Color == BLUE {
			b += .5
			bm += float64(p2.Intensity)/2.0
		}
	}
	var most int
	var m, mm float64
	if rm > gm && rm > bm {
		most = RED
		m = r
		mm = rm
	} else if gm > rm && gm > bm {
		most = GREEN
		m = g
		mm = gm
	} else if bm > gm && bm > rm {
		most = BLUE
		m = b
		mm = bm
	}
	if mm == 0 {
		return
	}

	if p.Color == WHITE {
		UnWhite(most, m, mm, p)
		return
	}
	if most == p.Color {
		var l, ll, w, ww float64
		if most == RED {
			l = g
			ll = gm
			w = b
			ww = bm
		} else if most == GREEN {
			l = b
			ll = bm
			w = r
			ww = rm
		} else {
			l = r
			ll = rm
			w = g
			ww = gm
		}
		if w == 0 && l == 0 {
			if m > 5 {
				Lose(p, WHITE)
			} else if m > 4 && mm/m > float64(p.Intensity) {
				Win(p)
			}
			return
		}
		if w >= l {
			Lose(p, WHITE)
		} else {
			Win(p)
		}
		ll += ww
		ww += ll
		return
		/*
		if m > 5 {
			Lose(p, WHITE)
		} else if m > 2 {
			Win(p)
		}
		return
		*/
	}
	if Beats(p.Color, most) {
		Win(p)
	} else {
		Lose(p, most)
	}
}

func UnWhite(color int, num, total float64, p *Point) {
	i := total/num/2
	if i < 1 {
		return
	}
	p.Color = color
	p.Intensity = int(i)
}

/*
func StepOne(x, y, width, height int, matrix *[][]Point) {
	dx, dy := Dir()
	x2 := x + dx
	y2 := y + dy
	p1 := &(*matrix)[x][y]
	if x2 < 0 || x2 >= width || y2 < 0 || y2 >= height {
		Lose(p1, WHITE)
		return
	}
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
*/

func NextMatrix(width, height int, matrix *[][]Point) *[][]Point {
	m := make([][]Point, width)
	for x := range m {
		m[x] = make([]Point, height)
	}
	for x := range *matrix {
		for y := range (*matrix)[x] {
			m[x][y] = (*matrix)[x][y]
			StepOne(x, y, width, height, &m[x][y], matrix)
		}
	}
	return &m
}

func makeMatrix(width, height int) *[][]Point {
	m := make([][]Point, width)
	for x := range m {
		m[x] = make([]Point, height)
		for y := range m[x] {
			m[x][y].Color = WHITE
			m[x][y].Intensity = 10
		}
	}

	/*
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
	*/

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
	for a := x; a < x+w; a++ {
		for b := y; b < y+h; b++ {
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
			gl.Vertex2f(gl.Float(x)*xby, gl.Float(y)*yby)
		}
	}
}
