package pkg

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

type FileCompareConfig struct {
	ReferenceApplication string        `json:"ReferenceApplication"`
	CompareApplications  []Application `json:"CompareApplications"`
	FilesToCompare       []string      `json:"FilesToCompare"`
}

type Application struct {
	ApplicationName string `json:"ApplicationName"`
	Path            string `json:"Path"`
}

func InitConfigDir() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dirname, ".fileCompare")); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(filepath.Join(dirname, ".fileCompare"), 0777)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func InitConfigFile() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dirname, ".fileCompare", ".fileCompare.json")); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(filepath.Join(dirname, ".fileCompare", ".fileCompare.json"))
		if err != nil {
			log.Fatal(err)
		}
		initJsonData(file)
	}
}

func initJsonData(file *os.File) {
	result := FileCompareConfig{}
	stringSlice := make([]string, 1)
	stringSlice[0] = "put data here"
	result.ReferenceApplication = "\\path\\to\\reference\\folder"
	result.CompareApplications = []Application{
		{
			ApplicationName: "App A",
			Path:            "\\path\\to\\application_A\\folder",
		},
	}
	result.FilesToCompare = []string{"compare.json"}
	byteValue, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	// Write back to file

	_, err = file.Write(byteValue)
	if err != nil {
		log.Fatal(err)
	}
}
