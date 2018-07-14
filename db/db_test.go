package db_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/PaluMacil/dwn/db"
)

var Db *db.Db

func TestMain(m *testing.M) {
	const testDir = "data-test"
	purgeDb(testDir)
	var err error
	Db, err = db.New(testDir)
	if err != nil {
		log.Fatalln("could not create db:", err)
	}
	retCode := m.Run()
	Db.Close()
	os.Exit(retCode)
}

func purgeDb(dir string) {
	dirRead, _ := os.Open(dir)
	dirFiles, _ := dirRead.Readdir(0)

	// Loop over the directory's files.
	for index := range dirFiles {
		fileHere := dirFiles[index]

		// Get name of file and its full path.
		nameHere := fileHere.Name()
		fullPath := filepath.Join(dir, nameHere)

		// Remove the file.
		err := os.Remove(fullPath)
		if err != nil {
			log.Println("While purging db files", err)
		}
	}
}
