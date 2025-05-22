package main

import (
	"encoding/csv"
	"log"
	"os"
)

func createCSVFile(filePath string, header []string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Unable to create CSV file:", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Unable to close CSV file:", err)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		log.Fatal("Unable to write record to CSV file:", err)
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func addCSVLine(data []string) {
	if !fileExists("./history.csv") {
		headers := []string{"Call time", "Number called", "Speech", "Can intervene"}
		createCSVFile("./history.csv", headers)
	}

	file, err := os.OpenFile("./history.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal("Unable to open CSV file:", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Unable to close CSV file:", err)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(data); err != nil {
		log.Fatal("Unable to write record to CSV file:", err)
	}
}
