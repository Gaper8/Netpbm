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
	max           uint8
}

// This function takes an image of type pgm as PGM parameter and returns PGM structure that represents the image.
func ReadPGM(filename string) (*PGM, error) {
	PGMthree := PGM{}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Le fichier ne s'ouvre pas.")
		return nil, err
	}
	defer file.Close()
	// I create PGMthree variable and assign the PGM structure to it. Then I open my file and show an error if the file does not open.

	scanner := bufio.NewScanner(file)

	booleathree := false
	booleafour := false
	booleafive := false
	lineone := 0

	// I create and initialize five variables

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
			// I scan my file and if the lines start with # it continues, to ignore them.
		} else if !booleathree {
			PGMthree.magicNumber = (scanner.Text())
			booleathree = true
			// Here, if there is no # it goes to this else if and we enter it if my booleathree variable is false. So I assign the line read to my PGMthree.magicNumber variable. And then I set my booleathree variable to true to no longer enter this condition.
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
			// Same here I enter this condition if booleafour is false. Then I take the line that the scanner reads and I separate the data, with the space. Then I convert the first data into an integer which I assign to the width variable. I do the same for height. Then I set booleafour to true to no longer enter this condition.

			PGMthree.data = make(([][]uint8), PGMthree.height)
			for i := range PGMthree.data {
				PGMthree.data[i] = make(([]uint8), PGMthree.width)
			}
			// Here I just initialize the matrix, my data array with the width and height that I retrieve above.
		} else if !booleafive {
			max, _ := strconv.Atoi(scanner.Text())
			if err != nil {
				return nil, err
			}
			PGMthree.max = uint8(max)
			booleafive = true
			// Here I also use an else if with booleafive false to enter the condition. Then I convert the read line and I assign this conversion to the max variable. Then I convert max to uint8 and assign it to PGMthree.max. Finally booleafive changes to true to no longer enter the condition.
		} else {

			if PGMthree.magicNumber == "P2" {
				pasdidee := strings.Fields(scanner.Text())
				for i := 0; i < PGMthree.width; i++ {
					value, _ := strconv.Atoi(pasdidee[i])
					PGMthree.data[lineone][i] = uint8(value)
				}
				lineone++
			}
			// Here I start by checking if the magicNumber is equal to P2, if this is the case I start by going through the line read and converting it into several strings, I store the result in the variable pasdidee. Then I scan the width. Then I convert each element of pasdidee to an integer and store the result in value. Finally I assign the converted values to the data matrix.
			if PGMthree.magicNumber == "P5" {
				databyte := scanner.Bytes()
				Indicedata := 0
				for y := 0; y < PGMthree.height; y++ {
					for x := 0; x < PGMthree.width; x++ {
						PGMthree.data[y][x] = databyte[Indicedata]
						Indicedata++
					}
				}
				break
			}
		}
	}
	// Here, if the magicnumber is not p2, it is p5 so we enter the condition. I initialize two variables, databyte contains bytes scanner which directly reads binary characters. I have two loops, one that runs the width and the other the height. Then I just assigned the databyte values to the location of PGMthree.data.

	return &PGM{PGMthree.data, PGMthree.width, PGMthree.height, PGMthree.magicNumber, PGMthree.max}, nil
	// I return PGMthree struct to my PGM pointer which contains all the image data.
}

func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// The size function has PGM pointer to the PGM structure, so to return the width and height, I just have to return the width and height of the PGM structure.

func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// For the At function I start by checking that my x and y coordinates are not outside my image. If the coordinates are included in the dimensions of my image then I return the coordinates of the pixel.

func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// The Set function allows you to modify the value of pixel. For this I give the location of the pixel and I change its value by assigning it the new value given as parameter with value.

func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)
	// For the save function, I start by opening my save file and checking if there is an error. Then I write the necessary information in the file.

	for _, pasdidee := range pgm.data {
		if pgm.magicNumber == "P2" {
			for _, pixel := range pasdidee {
				fmt.Fprintf(file, "%d ", pixel)
			}
			// Then I browse my data if the maqic number is P2 I write in my backup file each value as an integer, pixel being the variable or are the values to write.
			fmt.Fprintln(file)
		} else if pgm.magicNumber == "P5" {
			file.Write(pasdidee)
			// If the magic number is P5, I write all data to the file without conversion.
		}
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

// For Invert, I start by traversing the height and width of my image. Then I just inverted the value of each pixel at coordinates i and j. For this I change the current value to the opposite by making the max value minus the current value. This operation has the effect of giving the opposite value of the current pixel.

func (pgm *PGM) Flip() {
	var division int = (pgm.width / 2)
	var pasdidee uint8
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < division; j++ {
			pasdidee = pgm.data[i][j]
			pgm.data[i][j] = pgm.data[i][pgm.width-j-1]
			pgm.data[i][pgm.width-j-1] = pasdidee
		}
	}
}

// The Flip function is used to invert the image horizontally. To do this, I start by knowing how far I need to invert my data. This is the role of the division variable. Then all that remains is to swap the opposing points. For this we take the starting point and we invert it with the point at the same height but for the width we do width minus j (the value of our point) - 1 to be in agreement with the programming index which begins  0. Then we take the index of our opposite point which we replace with the opposite point which we store in pasdidee.

func (pgm *PGM) Flop() {
	var division int = (pgm.height / 2)
	var pasdidee uint8
	for i := 0; i < pgm.width; i++ {
		for j := 0; j < division; j++ {
			pasdidee = pgm.data[j][i]
			pgm.data[j][i] = pgm.data[pgm.height-j-1][i]
			pgm.data[pgm.height-j-1][i] = pasdidee
		}
	}
}

// The Flop function is used to invert the image vertically. To do this, I start by knowing how far I need to invert my data. This is the role of the division variable. Then all that remains is to swap the opposing points. For this we take the starting point and we invert it with the point at the same width but for the height we do height minus j (the value of our point) - 1 to be in agreement with the programming index which begins 0. Then we take the index of our opposite point which we replace with pasdidee, which is the starting point.

func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// This function allows me to change the value of the magic number by changing it to the value in the function parameter.

func (pgm *PGM) SetMaxValue(maxValue uint8) {
	if pgm.max != maxValue {
		pasdidee := float64(maxValue) / float64(pgm.max)
		pgm.max = maxValue
		for i := 0; i < len(pgm.data); i++ {
			for j := 0; j < len(pgm.data[i]); j++ {
				pgm.data[i][j] = uint8(float64(pgm.data[i][j]) * pasdidee)
			}
		}
	}
}

// This function allows me to change the value of max by changing it to maxValue in the function parameter. I first check if I don't have the same value for the max. Then, I calculate by what multiple to multiply each value. Then I go through my data and multiply each value by pasdidee.

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

// The RotateCW function allows you to rotate the image by degrees. To do this I start by creating a new matrix to be able to change the data between our two matrices. Then I iterate over the height and width. Then, I take pixel at its base position and so that it rotates 90 degrees the width becomes the height in the new data and vice versa but in addition for the height of datav2 I do calculation to have the new location of the value.

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

// The ToPBM function converts grayscale image to black and white image. I start by putting pointer to the pbm struct in variable and I put back the values that interest me. I also create new matrix. Then I iterate over the height and width. Finally, all I have to do is convert the gray pixels into black and white, so any pixel above half the max is white (true) otherwise black (false).
