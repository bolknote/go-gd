package main

import "gd"

func main() {
    pict := gd.CreateFromJpeg("IMG_1150.JPG")
    defer pict.Destroy()

    color := pict.ColorAllocate(0, 0, 0)

    // Non-Unicode font
    font := gd.GetFont(gd.FONTGIANT)

    pict.Char(font, 100, 100, "B", color)
    pict.String(font, 100, 120, "bolknote.ru", color)

    // Unicode font
    pict.StringFT(color, "/Library/Fonts/Impact.ttf", 12, 0, 100, 150, "Hello! Привет!")

    pict.Jpeg("out.jpg", 95)
}
