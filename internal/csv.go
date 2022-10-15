package internal

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/labstack/gommon/log"

	"github.com/dimchansky/utfbom"
)

func ReadCsv(filename string) ([][]string, error) {
	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		log.Print("Nothing to do.")
		os.Exit(0)
	}
	defer f.Close()

	// Read File into a Variable
	r := csv.NewReader(utfbom.SkipOnly(f))
	r.Comma = ';'
	lines, err := r.ReadAll()
	if err != nil {
		return [][]string{}, err
	}
	err = f.Close()
	if err != nil {
		log.Fatal("Closing file failed.")
	}
	return lines, nil
}

func createDir(path string) {
	os.MkdirAll(path, os.ModePerm)
}

func createFile(filename string, content string) {
	if !fileExists(filepath.Dir(filename)) {
		createDir(filepath.Dir(filename))
	}
	err := ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Creating file %s failed: %s", filename, err)
	}
	log.Infof("Created file %s.", filename)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
