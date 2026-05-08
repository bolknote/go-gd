package main

import (
	"fmt"

	gd "github.com/bolknote/go-gd/v2/pkg/gd"
)

func main() {
	img, err := gd.NewTrueColor(320, 180)
	must(err)
	defer func() { _ = img.Close() }()

	white, err := img.AllocateColor(255, 255, 255)
	must(err)
	black, err := img.AllocateColor(0, 0, 0)
	must(err)
	blue, err := img.AllocateColor(30, 120, 220)
	must(err)

	must(img.Fill(0, 0, white))
	must(img.FilledRectangle(20, 20, 300, 160, blue))
	must(img.FilledEllipse(160, 90, 120, 80, white))
	must(img.Line(20, 20, 300, 160, black))

	font, err := gd.BuiltinFont(gd.FontGiant)
	must(err)
	must(img.String(font, 100, 82, "go-gd v2", black))

	must(img.EncodePNGFile("out.png", nil))
	fmt.Printf("wrote out.png with libgd %s\n", gd.Version().String)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
