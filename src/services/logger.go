package services

import (
	"io"
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "", log.LstdFlags)

func SetOutput(w io.Writer) {
	Logger.SetOutput(w)
}
