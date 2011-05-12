package main

import "gd"

func main() {
    pict := gd.CreateFromJpeg("IMG_1150.JPG1")
    defer pict.Destroy()

    //pict.ImageJpeg("out.jpg", 95)
}
