package main

import (
	"log"
	"os"
	"sync"
)

type FileManager struct {
	filename string
	fileLock sync.Mutex
	file     *os.File
}

func (fman *FileManager) Add(game string) {
	if fman.file == nil {
		f, err := os.OpenFile(fman.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Printf("Couldn't open file: %s", err)
			return
		}
		fman.file = f
	}
	fman.fileLock.Lock()
	defer fman.fileLock.Unlock()
	fman.file.Write([]byte(game))
}
