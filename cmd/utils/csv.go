package utils

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
)

type CsvCompareWriter struct {
	writer            *csv.Writer
	headers           []string
	referenceFileName string
	fields            map[string][]string
}

func (csvCompare *CsvCompareWriter) InitCsv(path, referenceFileName string, compareFiles []string) func() {
	file, err := os.Create(path + ".csv")
	if err != nil {
		log.Fatalf("Error creating CSV file: %v", err)
	}
	if err := os.Chmod(path+".csv", 0777); err != nil {
		log.Fatalf("Error setting file permissions: %v\n", err)
	}
	csvCompare.headers = []string{"field", referenceFileName}
	csvCompare.referenceFileName = referenceFileName
	csvCompare.fields = make(map[string][]string)
	for _, compareFile := range compareFiles {
		baseFileName := filepath.Base(compareFile)
		csvCompare.headers = append(csvCompare.headers, baseFileName)
	}

	csvCompare.writer = csv.NewWriter(file)
	if err := csvCompare.writer.Write([]string{"", "Reference File"}); err != nil {
		log.Fatalf("Error writing row to CSV: %v", err)
	}

	if err := csvCompare.writer.Write(csvCompare.headers); err != nil {
		log.Fatalf("Error writing row to CSV: %v", err)
	}
	return func() {
		for _, field := range csvCompare.fields {
			err := csvCompare.writer.Write(field)
			if err != nil {
				log.Fatalf("Error writing row to CSV: %v", err)
			}
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatalf("Couldn't close file - %s - %s", path, err)
			}
		}(file)
		defer csvCompare.writer.Flush()
	}
}

func (csvCompare *CsvCompareWriter) InitCsvWithAppNames(path string, compareFiles []string, appNames []string) func() {
	file, err := os.Create(path + ".csv")
	if err != nil {
		log.Fatalf("Error creating CSV file: %v", err)
	}
	if err := os.Chmod(path+".csv", 0777); err != nil {
		log.Fatalf("Error setting file permissions: %v\n", err)
	}
	csvCompare.headers = []string{"field"}
	csvCompare.fields = make(map[string][]string)
	for _, compareFile := range compareFiles {
		baseFileName := filepath.Base(compareFile)
		csvCompare.headers = append(csvCompare.headers, baseFileName)
	}

	csvCompare.writer = csv.NewWriter(file)
	headers := []string{"", "Reference File"}
	headers = append(headers, appNames...)
	if err := csvCompare.writer.Write(headers); err != nil {
		log.Fatalf("Error writing row to CSV: %v", err)
	}

	if err := csvCompare.writer.Write(csvCompare.headers); err != nil {
		log.Fatalf("Error writing row to CSV: %v", err)
	}
	return func() {
		for _, field := range csvCompare.fields {
			err := csvCompare.writer.Write(field)
			if err != nil {
				log.Fatalf("Error writing row to CSV: %v", err)
			}
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatalf("Couldn't close file - %s - %s", path, err)
			}
		}(file)
		defer csvCompare.writer.Flush()
	}
}

func (csvCompare *CsvCompareWriter) WriteRow(referenceValue, fileValue interface{}, field, compareFileName string) {
	if len(csvCompare.fields[field]) == 0 {
		csvCompare.fields[field] = append(csvCompare.fields[field], field, ConvertToString(referenceValue))
	}

	csvCompare.fields[field] = append(csvCompare.fields[field], ConvertToString(fileValue))
}
