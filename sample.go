package main

import "gd"

func main() {
    pict := gd.CreateFromJpeg("IMG_1150.JPG")
    defer pict.Destroy()

    color := pict.ColorAllocate(0, 0, 0)
    font := gd.GetFont(gd.FONTLARGE)

    pict.Char(font, 100, 100, "B", color)

    pict.Jpeg("out.jpg", 95)
}
