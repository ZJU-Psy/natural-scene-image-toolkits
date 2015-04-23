/*
Project shoes
rainy @ 2015-04-23 <me@rainy.im>

1. Crop original image(s) according to given WIDTH/HEIGHT ratio;
2. Resize croped image(s) to the size WIDTH X HEIGHT;
3. Save to ./result.
*/
package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

var (
	WIDTH, HEIGHT      *int
	DESTPATH, FILENAME *string
)

func crop(pathsrc, pathdest string) image.Image {
	file, err := os.Open(pathsrc)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	ow, oh := img.Bounds().Max.X, img.Bounds().Max.Y
	if ow/oh > *WIDTH / *HEIGHT {
		ow = int(float64(*WIDTH) / float64(*HEIGHT) * float64(oh))
	} else {
		oh = int(float64(*HEIGHT) / float64(*WIDTH) * float64(ow))
	}

	croped, _ := cutter.Crop(img, cutter.Config{
		Width:   ow,
		Height:  oh,
		Options: cutter.Copy,
	})

	if pathdest != "" {
		out, err := os.Create(pathdest)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		// write new image to file
		jpeg.Encode(out, croped, nil)
	}

	return croped
}
func scale(img image.Image, dest string) error {
	m := resize.Resize(uint(*WIDTH), uint(*HEIGHT), img, resize.Lanczos3)
	out, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	return nil
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		fmt.Println("Resizing " + path + "...")
		croped := crop(path, "")
		os.Mkdir("./result", 0777)
		scale(croped, "./result/"+f.Name())
	}
	return nil
}
func main() {

	WIDTH = flag.Int("w", 400, "Resize width to [400]")
	HEIGHT = flag.Int("h", 400, "Resize height to [400]")
	FILENAME = flag.String("f", "", "Resize only this image")
	DESTPATH = flag.String("d", "", "Reisze all images under this dir")

	flag.Parse()

	if *FILENAME != "" {
		croped := crop(*FILENAME, "")
		scale(croped, "test.jpg")
	} else if *DESTPATH != "" {
		filepath.Walk("./test", visit)
	} else {
		fmt.Println("No file or dir given.")
		flag.Usage()
	}
}
