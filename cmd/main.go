package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sprakhar77/filereader/internal/reader"
)

func main() {
	filePathPtr := flag.String("f", "", "Enter the file path")
	datePtr := flag.String("d", "", "Enter the date")

	flag.Parse()
	if err := validateParams(filePathPtr, datePtr); err != nil {
		panic(err)
	}

	clr := reader.NewLogReader(*filePathPtr)

	if err := clr.Read(5 * reader.GB, 1); err != nil {
		panic(err)
	}

	cookies := clr.MostActiveCookies(*datePtr)
	for _, c := range cookies {
		fmt.Println(c)
	}
}

// validateParams validates the input flags
func validateParams(filePath, datePtr *string) error {
	if filePath == nil || len(*filePath) == 0 {
		return fmt.Errorf("file path cannot be empty")
	}

	if _, err := os.Stat(*filePath); err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	if datePtr == nil || len(*datePtr) == 0 {
		return fmt.Errorf("date cannot be empty")
	}

	_, err := time.Parse("2006-01-02", *datePtr)

	return err
}
