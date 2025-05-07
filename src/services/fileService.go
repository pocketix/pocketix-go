package services

import (
	"io"
	"log"
	"os"
)

func OpenFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	return data
}
