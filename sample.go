package main

import "gd"

func main() {
    pict := gd.CreateFromJpeg("IMG_1150.JPG")
    defer pict.Destroy()

    pict.Jpeg("out.jpg", 95)
}
