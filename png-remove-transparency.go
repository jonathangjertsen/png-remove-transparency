package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

type Flags struct {
	sources     []string
	destination string
}

func parseFlags(args []string) (*Flags, error) {
	flags := Flags{}
	if len(args) < 2 {
		return nil, fmt.Errorf("Need at least 2 (<source1> <source2> ... <sourceN> <destination>)\n")
	}
	flags.sources = args[:len(args)-1]
	flags.destination = args[len(args)-1]
	return &flags, nil
}

func fatal(err error) {
	fmt.Printf("%s\n", err)
	os.Exit(1)
}

func main() {
	args := os.Args[1:]

	// Parse inputs
	flags, err := parseFlags(args)
	if err != nil {
		fatal(err)
	}

	for _, source_file := range flags.sources {
		dest_file := filepath.Join(flags.destination, filepath.Base(source_file))
		fmt.Printf("%s -> %s\n", source_file, dest_file)

		// Open source file
		source, err := os.Open(source_file)
		if err != nil {
			fatal(err)
		}
		defer source.Close()

		// Decode source
		img_src, err := png.Decode(source)
		if err != nil {
			fatal(err)
		}

		// Update pixels
		bounds := img_src.Bounds()
		img_dest := image.NewNRGBA(bounds)
		x_min := bounds.Min.X
		x_max := bounds.Max.X
		y_min := bounds.Min.Y
		y_max := bounds.Max.Y
		for y := y_min; y < y_max; y++ {
			for x := x_min; x < x_max; x++ {
				c := img_src.At(x, y)
				r, g, b, a := c.RGBA()
				r_f := float32(r) / 65535 // why are these 16 bits
				g_f := float32(g) / 65535
				b_f := float32(b) / 65535
				a_f := float32(a) / 65535
				r_f += (1 - r_f) * (1 - a_f)
				g_f += (1 - g_f) * (1 - a_f)
				b_f += (1 - b_f) * (1 - a_f)
				img_dest.Set(x, y, color.RGBA{
					R: uint8(r_f * 255),
					G: uint8(g_f * 255),
					B: uint8(b_f * 255),
					A: 255,
				})
			}
		}

		// Open destination file
		dest, err := os.Create(dest_file)
		if err != nil {
			fatal(err)
		}
		defer dest.Close()

		// Write new file to it
		err = png.Encode(dest, img_dest)
		if err != nil {
			fatal(err)
		}
	}
}
