package main

import (
	"Netpbm"
)

func main() {

	pgm, _ := Netpbm.ReadPGM("testP5.pgm")

	pgm.SetMagicNumber("P5")

	pgm.Save("outputpgm.pgm")

}
