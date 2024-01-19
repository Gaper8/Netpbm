package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// This function takes an image of type pbm as a parameter and returns a structure that represents the image.
func ReadPBM(filename string) (*PBM, error) {
	PBMtwo := PBM{}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture")
		return nil, err
	}
	defer file.Close()
	// I create a variable and assign the PBM structure to it. Then I open my file and show an error if the file does not open.

	scanner := bufio.NewScanner(file)

	booleaone := false
	booleatwo := false
	line := 0
	// I create and initialize four variables

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
			// I scan my file and if the lines start with a # it continues, to ignore them.
		} else if !booleaone {
			PBMtwo.magicNumber = (scanner.Text())
			booleaone = true
			// Here, if there is no # it goes to this else if and we enter it if my booleaone variable is false. So I assign the line read to my PBMtwo.magicNumber variable. And then I set my booleaone variable to true to no longer enter this condition.
		} else if !booleatwo {
			size := strings.Split(scanner.Text(), " ")
			PBMtwo.width, err = strconv.Atoi(size[0])
			if err != nil {
				return nil, err
			}
			PBMtwo.height, err = strconv.Atoi(size[1])
			if err != nil {
				return nil, err
			}
			booleatwo = true
			// Same here I enter this condition if booleatwo is false. Then I take the line that the scanner reads and I separate the data, with the space. Then I convert the first data into an integer which I assign to the width variable. I do the same for height. Then I set booleatwo to true to no longer enter this condition.

			PBMtwo.data = make(([][]bool), PBMtwo.height)
			for i := range PBMtwo.data {
				PBMtwo.data[i] = make(([]bool), PBMtwo.width)
			}
			// Here I just initialize the matrix, my data array with the width and height that I retrieve above.
		} else {

			if PBMtwo.magicNumber == "P1" {
				a := strings.Fields(scanner.Text())
				for i := 0; i < PBMtwo.width; i++ {
					if a[i] == "1" {
						PBMtwo.data[line][i] = true
					} else {
						PBMtwo.data[line][i] = false
					}
				}
				line++
			}
			// Here I make a condition to process images in p1 format. I start by separating each value. Then I loop across the width, if the element read by the scanner is 1, true is stored at this location in the data and false otherwise. I ended up incrementing my line variable.
			if PBMtwo.magicNumber == "P4" {
				bytenumber := 0

				if PBMtwo.width%8 == 0 {
					bytenumber = (PBMtwo.width / 8)
				} else {
					bytenumber = (PBMtwo.width / 8) + 1
				}
				//Here is the second condition for the p4 format. I start by checking how many databyte I need to store a line of my image.
				padding := (bytenumber * 8) - PBMtwo.width
				binary := ToBinary(scanner.Text(), bytenumber, padding)
				colonne := 0
				linev3 := 0
				// Afterwards I check how many zeros I will have to add to complete the last byte. Then i initialize three variables. In binary I call a function that converts the binary line that the scanner reads into "readable" binary.
				for _, linev2 := range binary {
					if colonne == PBMtwo.width {
						linev3++
						colonne = 0
					}
					if linev2 == '1' {
						PBMtwo.data[linev3][colonne] = true
					} else {
						PBMtwo.data[linev3][colonne] = false
					}
					colonne++
				}
				// Then I store each converted element that is in binary. If my column variable is equal to my width, then that moves to the next a, and I reset column to zero for the rest. Finally when there is a 1, true is stored in my data and otherwise false.
			}
		}
	}
	fmt.Printf("%+v\n", PBMtwo)
	return &PBM{PBMtwo.data, PBMtwo.width, PBMtwo.height, PBMtwo.magicNumber}, nil
	// I return a pointer to my PBM struct which contains all the image data.
}

func ToBinary(test string, bytenumber, padding int) string {
	var result string = ""

	for i := 0; i < len(test); i++ {
		test := fmt.Sprintf("%08b", test[i])
		if i != 0 && (i+1)%bytenumber == 0 {
			test = test[:len(test)-padding]
		}
		result += test
	}
	return result
}

// I make a loop that goes through each element of my test variable. My test variable contains the conversion of each character into 8 databyte of binary. Then I check which byte needs padding, for that I check which byte is multiple of bytenumber. I finish by modifying my test variable so that it no longer contains the padding.

func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// The size function has a pointer to the PBM structure, so to return the width and height, I just have to return the width and height of the PBM structure.

func (pbm *PBM) At(x, y int) bool {
	if x < 0 || x > pbm.width || y < 0 || y > pbm.height {
		return false
	}
	return pbm.data[y][x]
}

// For the At function I start by checking that my x and y coordinates are not outside my image. If the coordinates are included in the dimensions of my image then I return the coordinates of the pixel.

func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// The Set function allows you to modify the value of a pixel. For this I give the location of the pixel and I change its value by assigning it the new value given as a parameter with value.

func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)
	// For the save function, I start by opening my save file and checking if there is an error. Then I write the necessary information in the file.

	if pbm.magicNumber == "P1" {
		for _, a := range pbm.data {
			for _, pixel := range a {
				if pixel {
					fmt.Fprint(file, "1 ")
				} else {
					fmt.Fprint(file, "0 ")
				}
			}
			fmt.Fprintln(file)
		}
		// If the magic number is P1, I browse my data. If in my data there is a true, a 1 is written otherwise a 0.
	} else if pbm.magicNumber == "P4" {

		for _, a := range pbm.data {
			databyte := make([]byte, (pbm.width+7)/8)
			for i, pixel := range a {
				if pixel {
					byte1 := i % 8
					byte2 := i / 8
					databyte[byte2] |= 1 << uint(7-byte1)
				}
			}
			file.Write(databyte)
		}
	}
	// If the magic number is P4, I browse my data. I start by browsing my data and I also create a byte array to store the compact binary values of the P4 format. Then I start by calculating which bit will be modified and write to the file. The next line is used to know which byte the bit will go into. Finally I modify the bytes in the table.
	return nil
}

func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

// For Invert, I start by iterating over the height and width of my image. Then I just inverted the value of each pixel at coordinates i and j. It's the ! which has the role of doing the inversion.

func (pbm *PBM) Flip() {
	var division int = (pbm.width / 2)
	var a bool
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < division; j++ {
			a = pbm.data[i][j]
			pbm.data[i][j] = pbm.data[i][pbm.width-j-1]
			pbm.data[i][pbm.width-j-1] = a
		}
	}
}

// The Flip function is used to invert the image horizontally. To do this, I start by knowing how far I need to invert my data. This is the role of the division variable. Then all that remains is to swap the opposing points. For this we take the starting point and we invert it with the point at the same height but for the width we do width minus j (the value of our point) - 1 to be in agreement with the programming index which begins a 0. Then we take the index of our opposite point which we replace with a, which is the starting point.

func (pbm *PBM) Flop() {
	var division int = (pbm.height / 2)
	var a bool
	for i := 0; i < pbm.width; i++ {
		for j := 0; j < division; j++ {
			a = pbm.data[j][i]
			pbm.data[j][i] = pbm.data[pbm.height-j-1][i]
			pbm.data[pbm.height-j-1][i] = a
		}
	}
}

// The Flop function is used to invert the image vertically. To do this, I start by knowing how far I need to invert my data. This is the role of the division variable. Then all that remains is to swap the opposing points. For this we take the starting point and we invert it with the point at the same width but for the height we do height minus j (the value of our point) - 1 to be in agreement with the programming index which begins a 0. Then we take the index of our opposite point which we replace with a, which is the starting point.

func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

// This function allows me to change the value of the magic number by changing it to the value in the function parameter.
