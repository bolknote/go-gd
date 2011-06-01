package main

import "gd"
import "fmt"

func main() {
    // http://www.php.net/manual/en/function.imagecreatefromjpeg.php
    pict := gd.CreateFromJpeg("1.jpg")

    // http://www.php.net/manual/en/function.imagedestroy.php
    defer pict.Destroy()

    pict.Sharpen(10)
    pict.Brightness(50)

    // http://www.php.net/manual/en/function.imagecolorallocate.php
    black := pict.ColorAllocate(0, 0, 0)
    white := pict.ColorAllocate(255, 255, 255)

    pict.SmoothEllipse(20, 20, 32, 32, white)

    // http://www.php.net/manual/en/function.imagefilledellipse.php
    pict.FilledEllipse(100, 100, 40, 50, white)

    // http://www.php.net/manual/en/function.imagecopyresampled.php
    pict.CopyResampled(pict, 40, 40, pict.Sx()-41, pict.Sy()-41, 20, 20, 40, 40)

    // Non-Unicode font
    font := gd.GetFont(gd.FONTGIANT)

    // http://www.php.net/manual/en/function.imagechar.php
    pict.Char(font, 100, 100, "B", black)
    // http://www.php.net/manual/en/function.imagestring.php
    pict.String(font, 100, 120, "bolknote.ru", black)

    // Unicode font
    fonts := gd.GetFonts()
    fmt.Printf("Found %d X11 TTF font(s)\n", len(fonts))

    if l := len(fonts); l > 0 {
        pict.StringFT(black, fonts[l-1], 12, 0, 100, 150, "Hello! Привет!")
    }

    // http://www.php.net/Imagejpeg
    pict.Jpeg("out.jpg", 95)
}
