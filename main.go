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

// structs
type img struct {
	file fs.DirEntry
}

type imgs struct {
	Files []fs.DirEntry
}

// files handler
func firstGorou() {
	br := false
	for {
		if br {
			endMainCh <- "end of main goroutine"
			break
		}

		select {
		case imgs := <-firstCh:
			for _, file := range imgs.Files {
				fmt.Println(file.Name())
			}
			br = true
		}
	}
}

func secondGorou() {}

func thirdGorou() {}

func main() {
	go firstGorou()

	fs, err := os.ReadDir(FOLDER)
	if err != nil {
		log.Println(ERR_READDIR)
		log.Fatal(err)
	}

	firstCh <- imgs{Files: fs}

	fmt.Println(<-endMainCh)
}
