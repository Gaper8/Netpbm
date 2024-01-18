package Netpbm

import (
	"bufio"
	"fmt"
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
			}
		}
		if PPMfor.magicNumber == "P6" {

		}
	}

	fmt.Printf("%+v\n", PPMfor)
	return &PPM{PPMfor.data, PPMfor.width, PPMfor.height, PPMfor.magicNumber, PPMfor.max}, nil

}

func (ppm *PPM) Size() (int, int) {
	width, height := ppm.height, ppm.width
	fmt.Printf("Largeur: %d, Hauteur: %d\n", width, height)
	return width, height
}

func (ppm *PPM) At(x, y int) Pixel {
	if x < 0 || x > ppm.width || y < 0 || y > ppm.height {
	}
	return ppm.data[y][x]
}

func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

func (ppm *PPM) Save(filename string) error {
	file, _ := os.Create(filename)
	defer file.Close()

	fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)

	for _, i := range ppm.data {
		for _, j := range i {
			fmt.Fprintf(file, "%d %d %d ", j.R, j.G, j.B)
		}
		fmt.Fprintln(file)
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
	var division int = (ppm.width / 2)
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
	var division int = (ppm.height / 2)
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
	ppm.max = maxValue
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
			brightness := uint8(0.299*float64(ppm.data[y][x].R) + 0.587*float64(ppm.data[y][x].G) + 0.114*float64(ppm.data[y][x].B))
			pgm.data[y][x] = brightness
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

	pointX := -1
	if p1.X < p2.X {
		pointX = 1
	}

	pointY := -1
	if p1.Y < p2.Y {
		pointY = 1
	}

	err := X - Y

	for {

		ppm.Set(p1.X, p1.Y, color)

		err2 := 2 * err

		if err2 > -Y {
			err -= Y
			p1.X += pointX
		}

		if err2 < X {
			err += X
			p1.Y += pointY
		}

		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}
	}
}

func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
}

func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
}

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	// ...
}

func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	// ...
}

func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	// ...
}

func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	// ...
}

func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	// ...
}

func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// ...
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
