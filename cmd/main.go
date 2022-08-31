package main

import (
	"fmt"
	"image"
	"io/fs"
	"log"
	"os"
	"time"
)

const (
	FOLDER    = "../images/"
	RESFOLDER = "../resized_images/"
)

const (
	ERR_READDIR       = "Read dir error"
	ERR_OPENINGFILE   = "os Opening file error"
	ERR_DECODINGFILE  = "png Decoding file error"
	ERR_CREATING_FILE = "os Create file error"
	ERR_ENCODING_PIC  = "png Encode pic error"
)

// chanels
var (
	firstCh   = make(chan imgs)
	secondCh  = make(chan img)
	thirdCh   = make(chan img)
	fourthCh  = make(chan img)
	endCh     = make(chan img)
	endMainCh = make(chan string)
)
var start = time.Now()

type img struct {
	file      fs.DirEntry
	index     int64
	openedImg *os.File
	saveFile  *os.File
	decoImg   image.Image
	resImg    image.Image
	finished  bool
}

type imgs struct {
	Files []fs.DirEntry
	count int64
}

// files handler
func selectCases() {
	var (
		lenFs    = int64(-1)
		counterF int64
	)

	for {
		if lenFs == counterF {
			endMainCh <- "end of main goroutine"
			break
		}

		select {
		case imgs := <-firstCh:
			lenFs = imgs.count
			for i, file := range imgs.Files {
				img := img{file: file, index: int64(i)}
				go img.openingGorou()
			}
		case f := <-secondCh:
			go f.decodingImgGorou()
		case f := <-thirdCh:
			go f.createFileToSave()
		case f := <-fourthCh:
			go f.resizeAndEncodeImg()
		case f := <-endCh:
			elapsed := time.Since(start)
			f.finished = true
			counterF++
			fmt.Println(f.file.Name(), f.index, counterF)
			fmt.Println("Execution time:")
			fmt.Println(elapsed.Seconds())
			fmt.Println("-+-+-+-+-+-")

		}
	}
}

func run() {
	go selectCases()

	fs, err := os.ReadDir(FOLDER)
	if err != nil {
		log.Println(ERR_READDIR)
		log.Fatal(err)
	}

	firstCh <- imgs{Files: fs, count: int64(len(fs))}

	fmt.Println(<-endMainCh)
}

func main() {
	run()
}
