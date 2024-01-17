package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

func ReadPGM(filename string) (*PGM, error) {
	PGMthree := PGM{}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture")
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	booleathree := false
	booleafour := false
	booleafive := false
	lineone := 0

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		} else if !booleathree {
			PGMthree.magicNumber = (scanner.Text())
			booleathree = true
		} else if !booleafour {
			size := strings.Split(scanner.Text(), " ")
			PGMthree.width, err = strconv.Atoi(size[0])
			if err != nil {
				return nil, err
			}
			PGMthree.height, err = strconv.Atoi(size[1])
			if err != nil {
				return nil, err
			}
			booleafour = true

			PGMthree.data = make(([][]uint8), PGMthree.height)
			for i := range PGMthree.data {
				PGMthree.data[i] = make(([]uint8), PGMthree.width)
			}
		} else if !booleafive {
			PGMthree.max, _ = strconv.Atoi(scanner.Text())
			booleafive = true
		} else {

			if PGMthree.magicNumber == "P2" {
				a := strings.Fields(scanner.Text())
				for i := 0; i < PGMthree.width; i++ {
					value, _ := strconv.Atoi(a[i])
					PGMthree.data[lineone][i] = uint8(value)
				}
				lineone++
			}
		}
		if PGMthree.magicNumber == "P5" {

		}
	}

	fmt.Printf("%+v\n", PGMthree)
	return &PGM{PGMthree.data, PGMthree.width, PGMthree.height, PGMthree.magicNumber, PGMthree.max}, nil

}

func (pgm *PGM) Size() (int, int) {
	width, height := pgm.height, pgm.width
	fmt.Printf("Largeur: %d, Hauteur: %d\n", width, height)
	return width, height
}

func (pgm *PGM) At(x, y int) uint8 {
	if x < 0 || x > pgm.width || y < 0 || y > pgm.height {
	}
	return pgm.data[y][x]
}

func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

func (pgm *PGM) Save(filename string) error {
	file, _ := os.Create(filename)
	defer file.Close()

	fmt.Fprintf(file, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	for _, i := range pgm.data {
		for _, j := range i {
			fmt.Fprintf(file, "%d ", j)
		}
		fmt.Fprintln(file)
	}
	return nil
}

func (pgm *PGM) Invert() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

func (pgm *PGM) Flip() {
	var division int = (pgm.width / 2)
	var a uint8
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < division; j++ {
			a = pgm.data[i][j]
			pgm.data[i][j] = pgm.data[i][pgm.width-j-1]
			pgm.data[i][pgm.width-j-1] = a
		}
	}
}

func (pgm *PGM) Flop() {
	var division int = (pgm.height / 2)
	var a uint8
	for i := 0; i < pgm.width; i++ {
		for j := 0; j < division; j++ {
			a = pgm.data[j][i]
			pgm.data[j][i] = pgm.data[pgm.height-j-1][i]
			pgm.data[pgm.height-j-1][i] = a
		}
	}
}

func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = int(maxValue)
}

func (pgm *PGM) Rotate90CW() {
	datav2 := make([][]uint8, pgm.width)
	for i := range datav2 {
		datav2[i] = make([]uint8, pgm.height)
	}

	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			datav2[x][pgm.height-y-1] = pgm.data[y][x]
		}
	}

	pgm.width, pgm.height = pgm.height, pgm.width
	pgm.data = datav2
}

func (pgm *PGM) ToPBM() *PBM {
	pbm := &PBM{
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1",
	}

	pbm.data = make([][]bool, pgm.height)
	for i := range pbm.data {
		pbm.data[i] = make([]bool, pgm.width)
	}

	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			pbm.data[y][x] = pgm.data[y][x] > uint8(pgm.max/2)
		}
	}

	return pbm
}
