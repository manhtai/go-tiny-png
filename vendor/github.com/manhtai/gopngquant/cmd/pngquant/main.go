package main

import (
	"flag"
	"github.com/manhtai/gopngquant"
	"log"
)

func main() {
	infile := flag.String("i", "", "input filename")
	outfile := flag.String("o", "", "output filename")
	speed := flag.Int("s", 3, "speed (1 slowest, 10 fastest)")

	flag.Parse()

	err := pngquant.CompressPngFile(*infile, *outfile, *speed)
	if err != nil {
		log.Fatalf("Compress fail: %s", err.Error())
	}
}
