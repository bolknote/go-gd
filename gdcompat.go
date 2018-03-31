package gd

import "path/filepath"
import "strings"
import "io/ioutil"
import . "math"
import "runtime"

func abs(i int) int {
	if i < 0 {
		return -i
	}

	return i
}

// Stack Blur Algorithm by Mario Klingemann <mario@quasimondo.com>
// "Go" language port by Evgeny Stepanischev http://bolknote.ru
func (img *Image) StackBlur(radius int, keepalpha bool) {
	if radius < 1 {
		return
	}

	w, h := int(img.Sx()), int(img.Sy())
	wm, hm, wh, div := w-1, h-1, w*h, radius*2+1

	len := map[bool]int{true: 3, false: 4}[keepalpha]

	rgba := make([][]byte, len)
	for i := 0; i < len; i++ {
		rgba[i] = make([]byte, wh)
	}

	vmin := make([]int, max(w, h))

	var x, y, i, yp, yi, yw, stackpointer, stackstart, rbs int
	var sir *[4]byte

	divsum := (div + 1) >> 1
	divsum *= divsum

	dv := make([]byte, 256*divsum)

	for i = 0; i < 256*divsum; i++ {
		dv[i] = byte(i / divsum)
	}

	yw, yi = 0, 0
	stack := make([][4]byte, div)
	r1 := radius + 1

	for y = 0; y < h; y++ {
		sum := make([]int, len)
		insum := make([]int, len)
		outsum := make([]int, len)

		for i = -radius; i <= radius; i++ {
			coords := yi + min(wm, max(i, 0))
			yc := coords / w
			xc := coords % w

			p := img.ColorsForIndex(img.ColorAt(xc, yc))

			sir = &stack[i+radius]
			sir[0] = (byte)(p["red"])
			sir[1] = (byte)(p["green"])
			sir[2] = (byte)(p["blue"])
			sir[3] = (byte)(p["alpha"])

			rbs = r1 - abs(i)
			for i := 0; i < len; i++ {
				sum[i] += int(sir[i]) * rbs
			}

			if i > 0 {
				for i := 0; i < len; i++ {
					insum[i] += int(sir[i])
				}
			} else {
				for i := 0; i < len; i++ {
					outsum[i] += int(sir[i])
				}
			}
		}

		stackpointer = radius

		for x = 0; x < w; x++ {
			for i := 0; i < len; i++ {
				rgba[i][yi] = dv[sum[i]]
				sum[i] -= outsum[i]
			}

			stackstart = stackpointer - radius + div
			sir = &stack[stackstart%div]

			for i := 0; i < len; i++ {
				outsum[i] -= int(sir[i])
			}

			if y == 0 {
				vmin[x] = min(x+radius+1, wm)
			}

			coords := yw + vmin[x]
			yc := coords / w
			xc := coords % w

			p := img.ColorsForIndex(img.ColorAt(xc, yc))

			sir[0] = byte(p["red"])
			sir[1] = byte(p["green"])
			sir[2] = byte(p["blue"])
			sir[3] = byte(p["alpha"])

			for i := 0; i < len; i++ {
				insum[i] += int(sir[i])
				sum[i] += insum[i]
			}

			stackpointer = (stackpointer + 1) % div
			sir = &stack[stackpointer%div]

			for i := 0; i < len; i++ {
				outsum[i] += int(sir[i])
				insum[i] -= int(sir[i])
			}

			yi++
		}

		yw += w
	}

	for x = 0; x < w; x++ {
		sum := make([]int, len)
		insum := make([]int, len)
		outsum := make([]int, len)

		yp = -radius * w

		for i = -radius; i <= radius; i++ {
			yi = max(0, yp) + x

			sir = &stack[i+radius]

			for i := 0; i < len; i++ {
				sir[i] = rgba[i][yi]
			}
			rbs = r1 - abs(i)

			for i := 0; i < len; i++ {
				sum[i] += int(rgba[i][yi]) * rbs
			}

			if i > 0 {
				for i := 0; i < len; i++ {
					insum[i] += int(sir[i])
				}
			} else {
				for i := 0; i < len; i++ {
					outsum[i] += int(sir[i])
				}
			}

			if i < hm {
				yp += w
			}
		}

		yi = x

		stackpointer = radius

		for y = 0; y < h; y++ {
			var alpha int

			if keepalpha {
				alpha = img.ColorsForIndex(img.ColorAt(yi%w, yi/w))["alpha"]
			} else {
				alpha = int(dv[sum[3]])
			}

			newpxl := img.ColorAllocateAlpha(int(dv[sum[0]]), int(dv[sum[1]]), int(dv[sum[2]]), alpha)
			if newpxl == -1 {
				newpxl = img.ColorClosestAlpha(int(dv[sum[0]]), int(dv[sum[1]]), int(dv[sum[2]]), alpha)
			}

			img.SetPixel(yi%w, yi/w, newpxl)

			for i := 0; i < len; i++ {
				sum[i] -= outsum[i]
			}

			stackstart = stackpointer - radius + div
			sir = &stack[stackstart%div]

			for i := 0; i < len; i++ {
				outsum[i] -= int(sir[i])
			}

			if x == 0 {
				vmin[y] = min(y+r1, hm) * w
			}

			p := x + vmin[y]

			for i := 0; i < len; i++ {
				sir[i] = rgba[i][p]
				insum[i] += int(sir[i])
				sum[i] += insum[i]
			}

			stackpointer = (stackpointer + 1) % div
			sir = &stack[stackpointer]

			for i := 0; i < len; i++ {
				outsum[i] += int(sir[i])
				insum[i] -= int(sir[i])
			}

			yi += w
		}
	}
}

// Originally written from scratch by Ulrich Mierendorff, 06/2006
// Rewritten and improved 04/2007, 07/2007, optimized circle version 03/2008
// "Go" language port by Evgeny Stepanischev http://bolknote.ru

func smootharcsegment(p *Image, cx, cy, a, b, aaAngleX, aaAngleY float64, fillColor Color, start, stop, seg float64) {
	xStart := Abs(a * Cos(start))
	yStart := Abs(b * Sin(start))
	xStop := Abs(a * Cos(stop))
	yStop := Abs(b * Sin(stop))

	dxStart, dyStart, dxStop, dyStop := float64(0), float64(0), float64(0), float64(0)

	color := p.ColorsForIndex(fillColor)

	if xStart != 0 {
		dyStart = yStart / xStart
	}

	if xStop != 0 {
		dyStop = yStop / xStop
	}

	if yStart != 0 {
		dxStart = xStart / yStart
	}

	if yStop != 0 {
		dxStop = xStop / yStop
	}

	aaStartX := Abs(xStart) >= Abs(yStart)
	aaStopX := xStop >= yStop

	for x := float64(0); x < a; x++ {
		_y1 := dyStop * x
		_y2 := dyStart * x

		var error1, error2 float64

		if xStart > xStop {
			error1 = _y1 - float64(int(_y1))
			error2 = 1 - _y2 + float64(int(_y2))

			_y1 -= error1
			_y2 += error2
		} else {
			error1 = 1 - _y1 + float64(int(_y1))
			error2 = _y2 - float64(int(_y2))

			_y1 += error1
			_y2 -= error2
		}

		switch seg {
		case 0:
			fallthrough
		case 2:
			var y1, y2 float64
			i := seg

			if !(start > i*Pi/2 && x > xStart) {
				var xp, yp, xa, ya float64

				if i == 0 {
					xp, yp, xa, ya = 1, -1, 1, 0
				} else {
					xp, yp, xa, ya = -1, 1, 0, 1
				}

				if stop < (i+1)*Pi/2 && x <= xStop {
					alpha := int(127 - float64(127-color["alpha"])*error1)
					diffColor1 := p.ColorExactAlpha(color["red"], color["green"], color["blue"], alpha)

					y1 = _y1

					if aaStopX {
						xx := int(cx + xp*x + xa)
						yy := int(cy + yp*(y1+1) + ya)

						p.SetPixel(xx, yy, diffColor1)
					}
				} else {
					y := b * Sqrt(1-Pow(x, 2)/Pow(a, 2))
					error := y - float64(int(y))
					y = float64(int(y))

					alpha := int(127 - float64(127-color["alpha"])*error)
					diffColor := p.ColorExactAlpha(color["red"], color["green"], color["blue"], alpha)

					y1 = y
					if x < aaAngleX {
						xx := int(cx + xp*x + xa)
						yy := int(cy + yp*(y1+1) + ya)

						p.SetPixel(xx, yy, diffColor)
					}
				}

				if start > i*Pi/2 && x <= xStart {
					alpha := int(127 - float64(127-color["alpha"])*error2)
					diffColor2 := p.ColorExactAlpha(color["red"], color["green"], color["blue"], alpha)

					y2 = _y2
					if aaStartX {
						xx := int(cx + xp*x + xa)
						yy := int(cy + yp*(y2-1) + ya)

						p.SetPixel(xx, yy, diffColor2)
					}
				} else {
					y2 = 0
				}

				if y2 <= y1 {
					xx := int(cx + xp*x + xa)
					yy1 := int(cy + yp*y1 + ya)
					yy2 := int(cy + yp*y2 + ya)

					p.Line(xx, yy1, xx, yy2, fillColor)
				}
			}

		case 1:
			fallthrough
		case 3:
			var y1, y2 float64
			i := seg

			if !(stop < (i+1)*Pi/2 && x > xStop) {
				var xp, yp, xa, ya float64

				if i == 1 {
					xp, yp, xa, ya = -1, -1, 0, 0
				} else {
					xp, yp, xa, ya = 1, 1, 1, 1
				}

				if start > i*Pi/2 && x < xStart {
					alpha := int(127 - float64(127-color["alpha"])*error2)
					diffColor2 := p.ColorExactAlpha(color["red"], color["green"], color["blue"], alpha)

					y1 = _y2
					if aaStartX {
						xx := int(cx + xp*x + xa)
						yy := int(cy + yp*(y1+1) + ya)

						p.SetPixel(xx, yy, diffColor2)
					}
				} else {
					y := b * Sqrt(1-Pow(x, 2)/Pow(a, 2))
					error := y - float64(int(y))
					y = float64(int(y))

					alpha := int(127 - float64(127-color["alpha"])*error)
					diffColor := p.ColorExactAlpha(color["red"], color["green"], color["blue"], alpha)

					y1 = y
					if x < aaAngleX {
						xx := int(cx + xp*x + xa)
						yy := int(cy + yp*(y1+1) + ya)

						p.SetPixel(xx, yy, diffColor)
					}
				}

				if stop < (i+1)*Pi/2 && x <= xStop {
					alpha := int(127 - float64(127-color["alpha"])*error1)
					diffColor1 := p.ColorExactAlpha(color["red"], color["green"], color["blue"], alpha)

					y2 = _y1
					if aaStopX {
						xx := int(cx + xp*x + xa)
						yy := int(cy + yp*(y2-1) + ya)
						p.SetPixel(xx, yy, diffColor1)
					}
				} else {
					y2 = 0
				}

				if y2 <= y1 {
					xx := int(cx + xp*x + xa)
					yy1 := int(cy + yp*y1 + ya)
					yy2 := int(cy + yp*y2 + ya)

					p.Line(xx, yy1, xx, yy2, fillColor)
				}
			}
		} // switch
	} // for x

	for y := float64(0); y < float64(b); y++ {
		_x1 := dxStop * y
		_x2 := dxStart * y

		var error1, error2 float64

		if yStart > yStop {
			error1 = _x1 - float64(int(_x1))
			error2 = 1 - _x2 - float64(int(_x2))
			_x1 -= error1
			_x2 += error2
		} else {
			error1 = 1 - _x1 + float64(int(_x1))
			error2 = _x2 + float64(int(_x2))
			_x1 += error1
			_x2 -= error2
		}

		switch seg {
		case 0:
			fallthrough
		case 2:
			var x1, x2 float64
			i := seg

			if !(start > i*Pi/2 && y > yStop) {
				var xp, yp, xa, ya float64

				if i == 0 {
					xp, yp, xa, ya = 1, -1, 1, 0
				} else {
					xp, yp, xa, ya = -1, 1, 0, 1
				}

				if stop < (i+1)*Pi/2 && y <= yStop {
					alpha := int(127 - float64(127-color["alpha"])*error1)
					diffColor1 := p.ColorExactAlpha(color["red"], color["green"], color["blue"], alpha)

					x1 = _x1
					if !aaStopX {
						xx := int(cx + xp*(x1-1) + xa)
						yy := int(cy + yp*y + ya)

						p.SetPixel(xx, yy, diffColor1)
					}
				}

				if start > i*Pi/2 && y < yStart {
					alpha := int(127 - float64(127-color["alpha"])*error2)
					diffColor2 := p.ColorExactAlpha(color["red"], color["green"], color["blue"], alpha)

					x2 = _x2
					if !aaStartX {
						xx := int(cx + xp*(x2+1) + xa)
						yy := int(cy + yp*y + ya)

						p.SetPixel(xx, yy, diffColor2)
					}
				} else {
					x := a * Sqrt(1-Pow(y, 2)/Pow(b, 2))
					error := x - float64(int(x))
					x = float64(int(x))

					alpha := int(127 - float64(127-color["alpha"])*error)
					diffColor := p.ColorExactAlpha(color["red"], color["green"], color["blue"], alpha)

					x1 = x
					if y < aaAngleY && y <= yStop {
						xx := int(cx + xp*(x1+1) + xa)
						yy := int(cy + yp*y + ya)

						p.SetPixel(xx, yy, diffColor)
					}
				}
			}

		case 1:
			fallthrough
		case 3:
			var x1, x2 float64
			i := seg

			if !(stop < (i+1)*Pi/2 && y > yStart) {
				var xp, yp, xa, ya float64

				if i == 1 {
					xp, yp, xa, ya = -1, -1, 0, 0
				} else {
					xp, yp, xa, ya = 1, 1, 1, 1
				}

				if start > i*Pi/2 && y < yStart {
					alpha := int(127 - float64(127-color["alpha"])*error2)
					diffColor2 := p.ColorExactAlpha(color["red"], color["green"], color["blue"], alpha)

					x1 = _x2
					if !aaStartX {
						xx := int(cx + xp*(x1-1) + xa)
						yy := int(cy + yp*y + ya)

						p.SetPixel(xx, yy, diffColor2)
					}
				}

				if stop < (i+1)*Pi/2 && y <= yStop {
					alpha := int(127 - float64(127-color["alpha"])*error1)
					diffColor1 := p.ColorExactAlpha(color["red"], color["green"], color["blue"], alpha)

					x2 = _x1
					if !aaStopX {
						xx := int(cx + xp*(x2+1) + xa)
						yy := int(cy + yp*y + ya)

						p.SetPixel(xx, yy, diffColor1)
					}
				} else {
					x := a * Sqrt(1-Pow(y, 2)/Pow(b, 2))
					error := x - float64(int(x))
					x = float64(int(x))

					alpha := int(127 - float64(127-color["alpha"])*error)
					diffColor := p.ColorExactAlpha(color["red"], color["green"], color["blue"], alpha)

					x1 = x
					if y < aaAngleY && y < yStart {
						xx := int(cx + xp*(x1+1) + xa)
						yy := int(cy + yp*y + ya)

						p.SetPixel(xx, yy, diffColor)
					}
				}
			}
		} // switch
	} // for y
}

func round(f float64) float64 {
	if f-float64(int(f)) >= 0.5 {
		return Ceil(f)
	}

	return Floor(f)
}

// Parameters:
// cx      - Center of ellipse, X-coord
// cy      - Center of ellipse, Y-coord
// w       - Width of ellipse ($w >= 2)
// h       - Height of ellipse ($h >= 2 )
// color   - Color of ellipse as a four component array with RGBA
// start   - Starting angle of the arc, no limited range!
// stop    - Stop     angle of the arc, no limited range!
// start _can_ be greater than $stop!

func (p *Image) SmoothFilledArc(cx, cy, w, h int, color Color, start, stop float64) {
	for start < 0 {
		start += 2 * Pi
	}

	for stop < 0 {
		stop += 2 * Pi
	}

	for start > 2*Pi {
		start -= 2 * Pi
	}

	for stop > 2*Pi {
		stop -= 2 * Pi
	}

	if start > stop {
		p.SmoothFilledArc(cx, cy, w, h, color, start, 2*Pi)
		p.SmoothFilledArc(cx, cy, w, h, color, 0, stop)

		return
	}

	a := round(float64(w) / 2)
	b := round(float64(h) / 2)
	fcx := float64(cx)
	fcy := float64(cy)

	aaAngle := Atan((b * b) / (a * a) * Tan(0.25*Pi))
	aaAngleX := a * Cos(aaAngle)
	aaAngleY := b * Sin(aaAngle)

	a -= 0.5
	b -= 0.5

	for i := float64(0); i < 4; i++ {
		if start < (i+1)*Pi/2 {
			if start > i*Pi/2 {
				if stop > (i+1)*Pi/2 {
					smootharcsegment(p, fcx, fcy, a, b, aaAngleX, aaAngleY, color, start, (i+1)*Pi/2, i)
				} else {
					smootharcsegment(p, fcx, fcy, a, b, aaAngleX, aaAngleY, color, start, stop, i)
					break
				}
			} else {
				if stop > (i+1)*Pi/2 {
					smootharcsegment(p, fcx, fcy, a, b, aaAngleX, aaAngleY, color, i*Pi/2, (i+1)*Pi/2, i)
				} else {
					smootharcsegment(p, fcx, fcy, a, b, aaAngleX, aaAngleY, color, i*Pi/2, stop, i)
					break
				}
			}
		}
	}
}

func (p *Image) SmoothFilledEllipse(cx, cy, w, h int, color Color) {
	p.SmoothFilledArc(cx, cy, w, h, color, 0, 2*Pi)
}

func searchfonts(dir string) (out []string) {
	files, e := ioutil.ReadDir(dir)
	if e == nil {
		for _, file := range files {
			if name := file.Name(); file.IsDir() {
				entry := filepath.Join(dir, name)
				out = append(out, searchfonts(entry)...)
			} else {
				if ext := filepath.Ext(name); ext != "" {
					ext := strings.ToLower(ext[1:])
					whitelist := []string{"ttf", "otf", "cid", "cff", "pcf", "fnt", "bdr", "pfr", "pfa", "pfb", "afm"}

					for _, wext := range whitelist {
						if ext == wext {
							out = append(out, name)
							break
						}
					}
				}
			}
		}
	}

	return
}

func GetFonts() (list []string) {
	fontpath, pathseparator := "", ""

	switch runtime.GOOS {
	case "darwin":
		fontpath, pathseparator = "/usr/share/fonts/truetype:/System/Library/Fonts:/Library/Fonts", ":"

	case "windows":
		fontpath, pathseparator = `C:\WINDOWS\FONTS;C:\WINNT\FONTS`, ";"

	default:
		fontpath, pathseparator = "/usr/X11R6/lib/X11/fonts/TrueType:/usr/X11R6/lib/X11/fonts/truetype:"+
			"/usr/X11R6/lib/X11/fonts/TTF:/usr/share/fonts/TrueType:/usr/share/fonts/truetype:"+
			"/usr/openwin/lib/X11/fonts/TrueType:/usr/X11R6/lib/X11/fonts/Type1:/usr/lib/X11/fonts/Type1:"+
			"/usr/openwin/lib/X11/fonts/Type1", ":"
	}

	if fontpath != "" {
		for _, dir := range strings.Split(fontpath, pathseparator) {
			list = append(list, searchfonts(dir)...)
		}
	}

	return
}

func (p *Image) filter(flt func(r, g, b, a, x, y int) (int, int, int, int)) {
	sx, sy := p.Sx(), p.Sy()
	for y := 0; y < sy; y++ {
		for x := 0; x < sx; x++ {
			rgba := p.ColorsForIndex(p.ColorAt(x, y))
			r, g, b, a := flt(rgba["red"], rgba["green"], rgba["blue"], rgba["alpha"], x, y)

			newpxl := p.ColorAllocateAlpha(r, g, b, a)
			if newpxl == -1 {
				newpxl = p.ColorClosestAlpha(r, g, b, a)
			}

			p.SetPixel(x, y, newpxl)
		}
	}
}

func (p *Image) GrayScale() {
	p.filter(func(r, g, b, a, x, y int) (int, int, int, int) {
		c := (int)(.299*float64(r) + .587*float64(g) + .114*float64(b))

		return c, c, c, a
	})
}

func (p *Image) Negate() {
	p.filter(func(r, g, b, a, x, y int) (int, int, int, int) {
		r = 255 - r
		g = 255 - g
		b = 255 - b

		return r, g, b, a
	})
}

func min(n1, n2 int) int {
	if n1 < n2 {
		return n1
	}

	return n2
}

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}

	return n2
}

func (p *Image) Brightness(brightness int) {
	if brightness == 0 {
		return
	}

	p.filter(func(r, g, b, a, x, y int) (int, int, int, int) {
		r = min(255, max(r+brightness, 0))
		g = min(255, max(g+brightness, 0))
		b = min(255, max(b+brightness, 0))

		return r, g, b, a
	})
}

func (p *Image) Contrast(contrast float64) {
	contrast = (100.0 - contrast) / 100.0
	contrast *= contrast

	corr := func(c int, contrast float64) int {
		f := float64(c) / 255.0
		f -= .5
		f *= contrast
		f += .5
		f *= 255.0

		return min(255, max(0, int(f)))
	}

	p.filter(func(r, g, b, a, x, y int) (int, int, int, int) {
		r = corr(r, contrast)
		g = corr(g, contrast)
		b = corr(b, contrast)

		return r, g, b, a
	})
}

func (p *Image) Color(r, g, b, a int) {
	p.filter(func(ri, gi, bi, ai, x, y int) (int, int, int, int) {
		ri = max(0, min(255, r+ri))
		gi = max(0, min(255, g+gi))
		bi = max(0, min(255, b+bi))
		ai = max(0, min(255, a+ai))

		return ri, gi, bi, ai
	})
}

func (p *Image) Convolution(filter [3][3]float32, filter_div, offset float32) {
	sx, sy := p.Sx(), p.Sy()
	srcback := CreateTrueColor(sx, sy)
	defer srcback.Destroy()

	srcback.SaveAlpha(true)
	newpxl := srcback.ColorAllocateAlpha(0, 0, 0, 127)
	srcback.Fill(0, 0, newpxl)

	p.Copy(srcback, 0, 0, 0, 0, sx, sy)

	var af func(p *Image, c int) int

	if p.TrueColor() {
		af = func(p *Image, c int) int { return (c & 0x7F000000) >> 24 }
	} else {
		af = func(p *Image, c int) int { return (int)((*p.img).alpha[c]) }
	}

	for y := 0; y < sy; y++ {
		for x := 0; x < sx; x++ {
			newr, newg, newb := float32(0), float32(0), float32(0)

			for j := 0; j < 3; j++ {
				yv := min(max(y-1+j, 0), sy-1)

				for i := 0; i < 3; i++ {
					pxl := srcback.ColorsForIndex(srcback.ColorAt(min(max(x-1+i, 0), sx-1), yv))

					newr += float32(pxl["red"]) * filter[j][i]
					newg += float32(pxl["green"]) * filter[j][i]
					newb += float32(pxl["blue"]) * filter[j][i]
				}
			}

			newr = (newr / filter_div) + offset
			newg = (newg / filter_div) + offset
			newb = (newb / filter_div) + offset

			r := min(255, max(0, int(newr)))
			g := min(255, max(0, int(newg)))
			b := min(255, max(0, int(newb)))

			newa := af(srcback, int(srcback.ColorAt(x, y)))

			newpxl = p.ColorAllocateAlpha(r, g, b, newa)
			if newpxl == -1 {
				newpxl = p.ColorClosestAlpha(r, g, b, newa)
			}

			p.SetPixel(x, y, newpxl)
		}
	}
}

func (p *Image) GaussianBlur() {
	filter := [3][3]float32{{1.0, 2.0, 3.0}, {2.0, 4.0, 2.0}, {1.0, 2.0, 1.0}}
	p.Convolution(filter, 16, 0)
}

func (p *Image) EdgeDetectQuick() {
	filter := [3][3]float32{{-1.0, 0.0, -1.0}, {0.0, 4.0, 0.0}, {-1.0, 0.0, -1.0}}
	p.Convolution(filter, 1, 127)
}

func (p *Image) Emboss() {
	filter := [3][3]float32{{1.5, 0.0, 0.0}, {0.0, 0.0, 0.0}, {0.0, 0.0, -1.5}}
	p.Convolution(filter, 1, 127)
}

func (p *Image) MeanRemoval() {
	filter := [3][3]float32{{-1, -1, -1}, {-1, 9, -1}, {-1, -1, -1}}
	p.Convolution(filter, 1, 0)
}

func (p *Image) Smooth(weight float32) {
	filter := [3][3]float32{{1, 1, 1}, {1, weight, 1}, {1, 1, 1}}
	p.Convolution(filter, weight+8, 0)
}
