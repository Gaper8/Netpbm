package main

import (
	"Netpbm"
)

func main() {

	pgm, _ := Netpbm.ReadPBM("fichierp1.pbm")

	pgm.SetMagicNumber("P4")

	pgm.Save("output.pbm")

}
