package main

import (
	"os"
	"fmt"
	"github.com/gographics/imagick/imagick"
	"encoding/base64"
)

type XNGFrame struct {
	data []byte
	delay uint
}

func gif2xng(filename string, frames *[]XNGFrame) (width uint, height uint) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Clear()

	err := mw.ReadImage(filename)
	if err != nil {
		panic(err)
	}

	mw.SetFirstIterator()

	cmw := mw.CoalesceImages()
	defer cmw.Clear()

	width = cmw.GetImageWidth()
	height = cmw.GetImageHeight()
	frameNum := cmw.GetNumberImages()

	for i := 0; i < int(frameNum); i++ {
		frame := XNGFrame {
			data: cmw.GetImageBlob(),
			delay: cmw.GetImageDelay(),
		}
		*frames = append(*frames, frame)
		cmw.NextImage()
	}

	return width, height
}

func writeXNG(filename string, width uint, height uint, frames *[]XNGFrame) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	frameNum := len(*frames)

	fmt.Fprintf(file, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:A=\"http://www.w3.org/1999/xlink\" width=\"%d\" height=\"%d\">", width, height)

	for i := 0; i < int(frameNum); i++ {
		fmt.Fprintf(
			file,
			"<image id=\"%06d\" height=\"100%%\" A:href=\"data:image/jpeg;base64,%s\"/>",
			i,
			base64.StdEncoding.EncodeToString((*frames)[i].data),
		)
		fmt.Fprint(file, (*frames)[i].delay)
	}

	for i := 0; i < int(frameNum); i++ {
		var begin string
		if (i == 0) {
			begin = fmt.Sprintf("A%06d.end; 0s", frameNum - 1)
		} else {
			begin = fmt.Sprintf("A%06d.end", i - 1)
		}
		fmt.Fprintf(
			file,
			"<set A:href=\"#%06d\" id=\"A%06d\" attributeName=\"width\" to=\"100%%\" dur=\"%dms\" begin=\"%s\"/>",
			i,
			i,
			(*frames)[i].delay * 10,
			begin,
		)
	}

	fmt.Fprintf(file, "</svg>")
}

func main() {
	argsLen := len(os.Args)
	if argsLen < 3 {
		fmt.Println("Usage:")
		fmt.Println("    gif2xng infile outfile")
		os.Exit(1)
	}

	infilename := os.Args[1]
	outfilename := os.Args[2]

	var frames []XNGFrame
	width, height := gif2xng(infilename, &frames)
	writeXNG(outfilename, width, height, &frames)
}
