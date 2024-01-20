package Netpbm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

type Pixel struct {
	R, G, B uint8
}

func ReadPPM(filename string) (*PPM, error) {
	PPMfor := PPM{}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture")
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	booleasix := false
	booleaseven := false
	booleaheight := false
	lineone := 0

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		} else if !booleasix {
			PPMfor.magicNumber = (scanner.Text())
			booleasix = true
		} else if !booleaseven {
			size := strings.Split(scanner.Text(), " ")
			PPMfor.width, err = strconv.Atoi(size[0])
			if err != nil {
				return nil, err
			}
			PPMfor.height, err = strconv.Atoi(size[1])
			if err != nil {
				return nil, err
			}
			booleaseven = true

			PPMfor.data = make(([][]Pixel), PPMfor.height)
			for i := range PPMfor.data {
				PPMfor.data[i] = make(([]Pixel), PPMfor.width)
			}
		} else if !booleaheight {
			max, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return nil, err
			}
			PPMfor.max = uint8(max)
			booleaheight = true
		} else {

			if PPMfor.magicNumber == "P3" {
				a := strings.Fields(scanner.Text())
				for i := 0; i < PPMfor.width; i++ {
					r, _ := strconv.Atoi(a[i*3])
					g, _ := strconv.Atoi(a[i*3+1])
					b, _ := strconv.Atoi(a[i*3+2])
					PPMfor.data[lineone][i] = Pixel{uint8(r), uint8(g), uint8(b)}
				}
				lineone++

			} else if PPMfor.magicNumber == "P6" {
				datappmrgb := make([]byte, PPMfor.width*PPMfor.height*3)
				file, _ := os.ReadFile(filename)
				if err != nil {
					return nil, nil
				}
				copy(datappmrgb, file[len(file)-(PPMfor.width*PPMfor.height*3):])
				pixel := 0
				for y := 0; y < PPMfor.height; y++ {
					for x := 0; x < PPMfor.width; x++ {
						PPMfor.data[y][x].R = datappmrgb[pixel]
						PPMfor.data[y][x].G = datappmrgb[pixel+1]
						PPMfor.data[y][x].B = datappmrgb[pixel+2]
						pixel += 3
					}
				}
				break
			}
		}
	}

	fmt.Printf("%+v\n", PPMfor)
	return &PPM{PPMfor.data, PPMfor.width, PPMfor.height, PPMfor.magicNumber, PPMfor.max}, nil

}

func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)

	for _, riri := range ppm.data {
		if ppm.magicNumber == "P3" {
			for _, fifi := range riri {
				fmt.Fprintf(file, "%d %d %d ", fifi.R, fifi.G, fifi.B)
			}
			fmt.Fprintln(file)
		} else if ppm.magicNumber == "P6" {
			for _, loulou := range riri {
				file.Write([]byte{loulou.R, loulou.G, loulou.B})
			}
		}
	}

	return nil
}

func (ppm *PPM) Invert() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			ppm.data[i][j].R = uint8(ppm.max) - ppm.data[i][j].R
			ppm.data[i][j].G = uint8(ppm.max) - ppm.data[i][j].G
			ppm.data[i][j].B = uint8(ppm.max) - ppm.data[i][j].B
		}
	}
}

func (ppm *PPM) Flip() {
	division := (ppm.width / 2)
	var a Pixel
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < division; j++ {
			a = ppm.data[i][j]
			ppm.data[i][j] = ppm.data[i][ppm.width-j-1]
			ppm.data[i][ppm.width-j-1] = a
		}
	}
}

func (ppm *PPM) Flop() {
	division := (ppm.height / 2)
	var a Pixel
	for i := 0; i < ppm.width; i++ {
		for j := 0; j < division; j++ {
			a = ppm.data[j][i]
			ppm.data[j][i] = ppm.data[ppm.height-j-1][i]
			ppm.data[ppm.height-j-1][i] = a
		}
	}
}

func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

func (ppm *PPM) SetMaxValue(maxValue uint8) {
	if ppm.max != maxValue {
		pasdidee := float64(maxValue) / float64(ppm.max)
		ppm.max = maxValue
		for i := 0; i < len(ppm.data); i++ {
			for j := 0; j < len(ppm.data[i]); j++ {
				ppm.data[i][j].R = uint8(float64(ppm.data[i][j].R) * pasdidee)
				ppm.data[i][j].G = uint8(float64(ppm.data[i][j].G) * pasdidee)
				ppm.data[i][j].B = uint8(float64(ppm.data[i][j].B) * pasdidee)
			}
		}
	}
}

func (ppm *PPM) Rotate90CW() {
	datav2 := make([][]Pixel, ppm.width)
	for i := range datav2 {
		datav2[i] = make([]Pixel, ppm.height)
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			datav2[x][ppm.height-y-1] = ppm.data[y][x]
		}
	}

	ppm.width, ppm.height = ppm.height, ppm.width
	ppm.data = datav2
}

func (ppm *PPM) ToPGM() *PGM {
	pgm := &PGM{
		magicNumber: "P2",
		width:       ppm.width,
		height:      ppm.height,
		max:         ppm.max,
	}

	pgm.data = make([][]uint8, ppm.height)
	for i := range pgm.data {
		pgm.data[i] = make([]uint8, ppm.width)
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			gray := uint8((int(ppm.data[y][x].R) + int(ppm.data[y][x].G) + int(ppm.data[y][x].B)) / 3)
			pgm.data[y][x] = gray
		}
	}

	return pgm
}

func (ppm *PPM) ToPBM() *PBM {
	pbm := &PBM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1",
	}

	pbm.data = make([][]bool, ppm.height)
	for i := range pbm.data {
		pbm.data[i] = make([]bool, ppm.width)
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			brightness := (uint32(ppm.data[y][x].R/3) + uint32(ppm.data[y][x].G/3) + uint32(ppm.data[y][x].B/3))
			pbm.data[y][x] = brightness > uint32(ppm.max/2)
		}
	}

	return pbm
}

type Point struct {
	X, Y int
}

func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {

	X := p2.X - p1.X
	if X < 0 {
		X = -X
	}

	Y := p2.Y - p1.Y
	if Y < 0 {
		Y = -Y
	}

	Xv2 := -1
	if p1.X < p2.X {
		Xv2 = 1
	}

	Yv2 := -1
	if p1.Y < p2.Y {
		Yv2 = 1
	}

	err := X - Y

	for {
		if p1.X >= 0 && p1.X < ppm.width && p1.Y >= 0 && p1.Y < ppm.height {

			ppm.Set(p1.X, p1.Y, color)
		}

		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}

		err2 := 2 * err

		if err2 > -Y {
			err -= Y
			p1.X += Xv2
		}

		if err2 < X {
			err += X
			p1.Y += Yv2
		}

		if p1.X < 0 || p1.X >= ppm.width || p1.Y < 0 || p1.Y >= ppm.height {
			break
		}
	}
}
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {

	if width <= 0 || height <= 0 || p1.X < 0 || p1.Y < 0 || p1.X >= ppm.width || p1.Y >= ppm.height {
		return
	}
	if p1.X < 0 {
		p1.X = 0
	}
	if p1.Y < 0 {
		p1.Y = 0
	}

	if p1.X+width > ppm.width {
		width = ppm.width - p1.X
	}
	if p1.Y+height > ppm.height {
		height = ppm.height - p1.Y
	}

	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X + width, p1.Y + height}
	p4 := Point{p1.X, p1.Y + height}

	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p4, color)
	ppm.DrawLine(p4, p1, color)
}

func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	for i := 0; i < height; i++ {
		beginning := Point{p1.X, p1.Y + i}
		end := Point{p1.X + width, p1.Y + i}
		ppm.DrawLine(beginning, end, color)
	}
}

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {

	for x := 0; x < ppm.height; x++ {
		for y := 0; y < ppm.width; y++ {
			dx := float64(x) - float64(center.X)
			dy := float64(y) - float64(center.Y)
			distance := math.Sqrt(dx*dx + dy*dy)

			if math.Abs(distance-float64(radius)) < 1.0 && distance < float64(radius) {
				ppm.Set(x, y, color)
			}
		}
	}
	ppm.Set(center.X-(radius-1), center.Y, color)
	ppm.Set(center.X+(radius-1), center.Y, color)
	ppm.Set(center.X, center.Y+(radius-1), color)
	ppm.Set(center.X, center.Y-(radius-1), color)
}

func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	for radius >= 0 {
		ppm.DrawCircle(center, radius, color)
		radius--
	}
}

func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	for p1 != p2 {
		ppm.DrawLine(p3, p1, color)
		if p1.X != p2.X && p1.X < p2.X {
			p1.X++
		} else if p1.X != p2.X && p1.X > p2.X {
			p1.X--
		}
		if p1.Y != p2.Y && p1.Y < p2.Y {
			p1.Y++
		} else if p1.Y != p2.Y && p1.Y > p2.Y {
			p1.Y--
		}
	}
	ppm.DrawLine(p3, p1, color)
}

func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	sizepolygon := len(points)

	for i := 0; i < sizepolygon-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}
	ppm.DrawLine(points[sizepolygon-1], points[0], color)
}

func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	ppm.DrawPolygon(points, color)
	for i := 0; i < ppm.height; i++ {
		var placepixel []int
		var nombrepixels int
		for j := 0; j < ppm.width; j++ {
			if ppm.data[i][j] == color {
				nombrepixels += 1
				placepixel = append(placepixel, j)
			}
		}
		if nombrepixels > 1 {
			for a := placepixel[0] + 1; a < placepixel[len(placepixel)-1]; a++ {
				ppm.data[i][a] = color

			}
		}
	}
}

func (ppm *PPM) DrawKochSnowflake(n int, start Point, width int, color Pixel) {
	// N is the number of iterations.
	// Koch snowflake is a 3 times a Koch curve.
	// Start is the top point of the snowflake.
	// Width is the width all the lines.
	// Color is the color of the lines.
	// ...
}

func (ppm *PPM) DrawSierpinskiTriangle(n int, start Point, width int, color Pixel) {
	// N is the number of iterations.
	// Start is the top point of the triangle.
	// Width is the width all the lines.
	// Color is the color of the lines.
	// ...
}

func (ppm *PPM) DrawPerlinNoise(color1 Pixel, color2 Pixel) {
	// Color1 is the color of 0.
	// Color2 is the color of 1.
}

func (ppm *PPM) KNearestNeighbors(newWidth, newHeight int) {
	// ...
}
