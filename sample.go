package main

import "gd"

func main() {
    pict := gd.CreateFromJpeg("IMG_1150.JPG")
    defer pict.Destroy()

    color := pict.ColorAllocate(0, 0, 0)
    font := gd.GetFont(gd.FONTGIANT)

    pict.Char(font, 100, 100, "B", color)
    pict.String(font, 100, 120, "bolknote.ru", color)

    pict.Jpeg("out.jpg", 95)
}
