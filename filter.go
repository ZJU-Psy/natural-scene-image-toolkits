/*
PACKAGE main
rainy @ 2015-04-23 <me@rainy.im>

筛选图片，文件小于5K，尺寸高宽小于300将被剔除

Usage:

	-d ./test 	Destination path, './test' by default.
	-r 					Remove it directly, false by default.
*/
package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"os"
	"path/filepath"
)

var (
	DESTPATH *string
	REMOVEIT *bool
)

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		fmt.Printf("Visit: %s: %.0fK \t", f.Name(), float64(f.Size()*1.0/1024.0))

		// Open image file
		imgf, _ := os.Open(path)
		defer imgf.Close()
		img, _, _ := image.DecodeConfig(imgf)
		fmt.Printf("Image Size: [%dX%d]\n", img.Width, img.Height)

		if f.Size()*1.0/1024 < 5 || img.Width < 300 || img.Height < 300 {
			if *REMOVEIT {
				// Just remove it
				os.Remove(path)
				fmt.Printf("Remove %s...\n", f.Name())
			} else {
				fmt.Printf("Move %s to ./tmp ...\n", f.Name())
				os.MkdirAll("."+string(filepath.Separator)+"tmp", 0777)
				os.Rename(path, "."+string(filepath.Separator)+"tmp"+string(filepath.Separator)+f.Name())
			}
		}
	}

	return nil
}
func main() {
	DESTPATH := flag.String("d", "./test", "Dest dir - ./test by default.")
	REMOVEIT = flag.Bool("r", false, "Remove filtered images directly!")

	flag.Parse()
	fmt.Println(*DESTPATH)

	err := filepath.Walk(*DESTPATH, visit)
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}
