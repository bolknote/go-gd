package main

import "gd"
import "fmt"

func main() {
    pict := gd.CreateFromJpeg("IMG_1150.JPG")
    defer pict.Destroy()

    color := pict.ColorAllocate(0, 0, 0)

    // Non-Unicode font
    font := gd.GetFont(gd.FONTGIANT)

    pict.Char(font, 100, 100, "B", color)
    pict.String(font, 100, 120, "bolknote.ru", color)

    // Unicode font
    fonts := gd.GetTtfFonts()
    fmt.Printf("Found %d X11 TTF font(s)\n", len(fonts))

    if l := len(fonts); l > 0 {
        pict.StringFT(color, fonts[l-1], 12, 0, 100, 150, "Hello! Привет!")
    }

    pict.Jpeg("out.jpg", 95)
}
