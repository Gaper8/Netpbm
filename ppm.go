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

// This function takes an image of type pppm as PPM parameter and returns PPM structure that represents the image.
func ReadPPM(filename string) (*PPM, error) {
	PPMfor := PPM{}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Le fichier ne s'ouvre pas.")
		return nil, err
	}
	defer file.Close()
	// I create PGMfor variable and assign the PPM structure to it. Then I open my file and show an error if the file does not open.

	scanner := bufio.NewScanner(file)

	booleasix := false
	booleaseven := false
	booleaheight := false
	lineone := 0

	// I create and initialize five variables

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
			// I scan my file and if the lines start with # it continues, to ignore them.
		} else if !booleasix {
			PPMfor.magicNumber = (scanner.Text())
			booleasix = true
			// Here, if there is no # it goes to this else if and we enter it if my booleasix variable is false. So I assign the line read to my PGMfor.magicNumber variable. And then I set my booleasix variable to true to no longer enter this condition.
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
			// Same here I enter this condition if booleaseven is false. Then I take the line that the scanner reads and I separate the data, with the space. Then I convert the first data into an integer which I assign to the width variable. I do the same for height. Then I set booleaseven to true to no longer enter this condition.

			PPMfor.data = make(([][]Pixel), PPMfor.height)
			for i := range PPMfor.data {
				PPMfor.data[i] = make(([]Pixel), PPMfor.width)
			}
			// Here I just initialize the matrix, my data array with the width and height that I retrieve above.
		} else if !booleaheight {
			max, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return nil, err
			}
			PPMfor.max = uint8(max)
			booleaheight = true
			// Here I also use an else if with booleaheight false to enter the condition. Then I convert the read line and I assign this conversion to the max variable. Then I convert max to uint8 and assign it to PGMfor.max. Finally booleaheight changes to true to no longer enter the condition.
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
				// Here to enter this condition I start by checking if the magic number is p3, if this is the case I start by separating each character that the scanner reads. Then I travel my width. Then I make sure to know the index of each pixel, since they go from three to three. The calculation therefore helps me to know the first index of the first pixel and then we go from three to three to recover each value. Afterwards I just assign the data to the PPMfor.data variable.
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
	// Well here, if we don't fit into p3 we'll fit into p6. I start by creating a byte array of the size width multiplied by height multiplied by three for the rgb. Then I read the file and store it in file. Then I retrieve the binary data from the file in my datppmrgb variable. Then I iterate over my width and height. Finally I retrieve each pixel value from my datappmrgb data and I put them in the correct pixel of the PPMfor.data matrix. and of course I don't forget to increment pixel by three because a pixel has three values.

	return &PPM{PPMfor.data, PPMfor.width, PPMfor.height, PPMfor.magicNumber, PPMfor.max}, nil
	// I return PPMfor struct to my PPM pointer which contains all the image data.
}

func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// The size function has PPM pointer to the PPM structure, so to return the width and height, I just have to return the width and height of the PPM structure.

func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// For the At function I start by checking that my x and y coordinates are not outside my image. If the coordinates are included in the dimensions of my image then I return the coordinates of the pixel.

func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

// The Set function allows you to modify the value of pixel. For this I give the location of the pixel and I change its value by assigning it the new value given as parameter with value.

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

// I start by creating a save file. Then I enter the necessary information (height, width, magic number, max). Then, I browse my data, if it is p3 I enter the condition and I browse the pixels of the line and I enter the pixel values. If the magic number is p6, I scan the pixels again but this time the pixel values are written in binary format with file.Write.

func (ppm *PPM) Invert() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			ppm.data[i][j].R = uint8(ppm.max) - ppm.data[i][j].R
			ppm.data[i][j].G = uint8(ppm.max) - ppm.data[i][j].G
			ppm.data[i][j].B = uint8(ppm.max) - ppm.data[i][j].B
		}
	}
}

// I start by going through the height and width then to invert the values I assign the position i, j of my value to the opposite position by subtracting the amx from the current value. I do this for each component (r,g,b) of the pixel.

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

// The Flip function is used to invert the image horizontally. To do this, I start by knowing how far I need to invert my data. This is the role of the division variable. Then all that remains is to swap the opposing points. For this we take the starting point and we invert it with the point at the same height but for the width we do width minus j (the value of our point) - 1 to be in agreement with the programming index which begins  0. Then we take the index of our opposite point which we replace with the opposite point which we store in pasdidee.

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

// The Flop function is used to invert the image vertically. To do this, I start by knowing how far I need to invert my data. This is the role of the division variable. Then all that remains is to swap the opposing points. For this we take the starting point and we invert it with the point at the same width but for the height we do height minus j (the value of our point) - 1 to be in agreement with the programming index which begins 0. Then we take the index of our opposite point which we replace with pasdidee, which is the starting point.

func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// This function allows me to change the value of the magic number by changing it to the value in the function parameter.

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

// For setmaxvalue I start by checking that the value of the max from the beginning and the new one is not the same, then I calculate the value which will allow me to change the value of the components. To divide the new max with the old one I store in pasdidee, the float64 is used to have a precise calculation. Then I go through each pixel of the data matrix. Finally, to have the change in data value in accordance with the max I multiply the value at location i,j by pasdidee, and I convert to uint8.

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

// The RotateCW function allows you to rotate the image by degrees. To do this I start by creating a new matrix to be able to change the data between our two matrices. Then I iterate over the height and width. Then, I take pixel at its base position and so that it rotates 90 degrees the width becomes the height in the new data and vice versa but in addition for the height of datav2 I do calculation to have the new location of the value.

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

// I start by storing the PGM struct in my pgm variable. Then I create a new matrix. Then I go through each pixel of the image, and I calculate the average brightness of an entire pixel, for that I take each value I add them and I divide everything by 3. To finish I attribute the average gray value to the pixel of my data.

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

// I start by storing the PBM struct in my pbm variable. Then I create a new matrix. I scan each pixel in height and width. Then, I calculate the average value of a pixel by adding each value and dividing by 3. Then I assign my data the value true if the average of my pixel is greater than the max divided by 2.

type Point struct {
	X, Y int
}

// I create a Point structure with x and y ints.

func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {

	X := p2.X - p1.X
	if X < 0 {
		X = -X
	}

	Y := p2.Y - p1.Y
	if Y < 0 {
		Y = -Y
	}

	// These two conditions are used to calculate the margin of our line and also to know the direction in which we must go, thanks to the absolute value.

	Xv2 := -1 // In this case from right to left.
	if p1.X < p2.X {
		Xv2 = 1 // In this case, left to right.
	}
	// Then to be precise we want to know if our line should be drawn from right to left or vice versa.

	Yv2 := -1 // In this case from top to bottom.
	if p1.Y < p2.Y {
		Yv2 = 1 // In this case from bottom to top.
	}
	// And of course we also want to know here if the line will be drawn from top to bottom or from bottom to top.

	err := X - Y

	for {
		// Check that the pixel is within the image boundaries
		if p1.X >= 0 && p1.X < ppm.width && p1.Y >= 0 && p1.Y < ppm.height {

			ppm.Set(p1.X, p1.Y, color)
		}
		// Here we color the pixel

		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}
		// If the coloring is finished we break, since the loop no longer needs to be continued.

		err2 := 2 * err

		// This variable is used to determine when we need to move in the direction of Y

		if err2 > -Y {
			err -= Y
			p1.X += Xv2 // and we move forward
		}

		// If err2 is greater than the opposite of Y, we must move in the X direction

		if err2 < X {
			err += X
			p1.Y += Yv2 // and we move forward
		}

		// If err2 is greater than the opposite of X, we must move in the Y direction

		if p1.X < 0 || p1.X >= ppm.width || p1.Y < 0 || p1.Y >= ppm.height {
			break
		}
		// Finally, we check that the point is indeed within the limit of the image
	}
}
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {

	if width <= 0 || height <= 0 || p1.X < 0 || p1.Y < 0 || p1.X >= ppm.width || p1.Y >= ppm.height {
		return
	}

	// I check if the coordinates are outside the image boundaries.

	if p1.X < 0 {
		p1.X = 0
	}
	if p1.Y < 0 {
		p1.Y = 0
	}

	// Here I adjust the coordinates are negative. If so I adjust to zero.

	if p1.X+width > ppm.width {
		width = ppm.width - p1.X
	}
	if p1.Y+height > ppm.height {
		height = ppm.height - p1.Y
	}

	// Here I adjust the width and height. If this data is outside the limits of the image I reduce it to avoid overflow.

	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X + width, p1.Y + height}
	p4 := Point{p1.X, p1.Y + height}

	// Create the 3 corners of the rectangle

	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p4, color)
	ppm.DrawLine(p4, p1, color)

	// We link them all so as to make a loop
}

func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	for i := 0; i < height; i++ {
		beginning := Point{p1.X, p1.Y + i}
		end := Point{p1.X + width, p1.Y + i}
		ppm.DrawLine(beginning, end, color)
	}
}

//I start by scanning the height. I initialize beginning which stores the position of the start of the line (p1.X for the abscissa and p1.Y + 1 as the ordinate. And the variable end bas is the same principle but for the end. Finally, I have just call drawline to draw a line between beginnig and end.

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

// I start by going through each pixel of the image by traversing the height and width. Then I calculate the distance between the center of the circle and each pixel. Then, if the distance of a pixel is very close to the radius of the circle and the overall distance is less than the radius of the circle, then I color that pixel. Finally, I color four additional pixels, which are the edges of the circle.

func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	for radius >= 0 {
		ppm.DrawCircle(center, radius, color)
		radius--
	}
}

// Here I call my drawcircle function which draws a circle of the specified size then I decrement the radius so that with each iteration a smaller circle is drawn until the radius is no longer >= 0.

func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// Here it's very simple, I call drawline so that it draws lines between each point of the triangle as a parameter.

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

// To color the triangle what I do is that as long as p1 is not equal to the coordinates of p2, a line is drawn between p1 and p3. So then I make sure that p1 relates to each iteration of p2. This solution has the effect of filling the triangle line by line with p3 as a reference point.

func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	sizepolygon := len(points)

	for i := 0; i < sizepolygon-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}
	ppm.DrawLine(points[sizepolygon-1], points[0], color)
}

// I start by storing in my sizepolygon variable how many points the polygon contains. Then I create a line between all points except the last one. The last point I draw outside the loop.

func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	ppm.DrawPolygon(points, color)
	for i := 0; i < ppm.height; i++ {
		var placepixel []int
		var numberpixels int
		for j := 0; j < ppm.width; j++ {
			if ppm.data[i][j] == color {
				numberpixels += 1
				placepixel = append(placepixel, j)
			}
		}
		if numberpixels > 1 {
			for a := placepixel[0] + 1; a < placepixel[len(placepixel)-1]; a++ {
				ppm.data[i][a] = color

			}
		}
	}
}

// I start by going through each row of the image in height. Then, I check and store in placepixel all the pixels already placed on the lines (the outline of the polygon) and of the right color. Finally, I fill the pixels inside the polygon.

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
