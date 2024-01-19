package main

import (
	"Netpbm"
)

func main() {

	pgm, _ := Netpbm.ReadPGM("testP2.pgm")

	pgm.SetMagicNumber("P2")

	pgm.Save("outputpgm.pgm")

}
