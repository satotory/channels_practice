package main

import (
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
)

func (f img) openingGorou() {
	var err error

	f.openedImg, err = os.Open(FOLDER + f.file.Name())
	if err != nil {
		log.Println(ERR_OPENINGFILE)
		log.Fatal(err)
	}

	secondCh <- f
}

func (f img) createFileToSave() {
	var err error

	f.saveFile, err = os.Create(RESFOLDER + f.file.Name())
	if err != nil {
		log.Println(ERR_CREATING_FILE)
		log.Fatal(err)
	}

	fourthCh <- f
}

func (f img) decodingImgGorou() {
	var err error

	f.decoImg, err = jpeg.Decode(f.openedImg)
	if err != nil {
		log.Println(ERR_DECODINGFILE)
		log.Fatal(err)
	}
	f.openedImg.Close()

	thirdCh <- f
}

func (f img) resizeAndEncodeImg() {
	defer f.saveFile.Close()
	f.resImg = resize.Thumbnail(300, 200, f.decoImg, resize.Lanczos3)

	err := jpeg.Encode(f.saveFile, f.resImg, nil)
	if err != nil {
		log.Println()
		log.Fatal(err)
	}
	endCh <- f
}
