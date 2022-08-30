package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
)

const FOLDER = "./images"
const ERR_READDIR = "Read dir error"

// chanels
var (
	firstCh   = make(chan imgs)
	secondCh  = make(chan img)
	thirdCh   = make(chan img)
	endCh     = make(chan bool)
	endMainCh = make(chan string)
)

type img struct {
	file     fs.DirEntry
	index    int64
	finished bool
}

type imgs struct {
	Files []fs.DirEntry
	count int64
}

// files handler
func firstGorou() {
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
				go secondGorou(file, i)
			}
		case f := <-secondCh:
			f.finished = true
			counterF++
			fmt.Println(f.file.Name(), f.index, counterF)
		}
	}
}

func secondGorou(f fs.DirEntry, i int) {
	secondCh <- img{file: f, index: int64(i)}
}

func thirdGorou() {}

func main() {
	go firstGorou()

	fs, err := os.ReadDir(FOLDER)
	if err != nil {
		log.Println(ERR_READDIR)
		log.Fatal(err)
	}

	firstCh <- imgs{Files: fs, count: int64(len(fs))}

	fmt.Println(<-endMainCh)
}
