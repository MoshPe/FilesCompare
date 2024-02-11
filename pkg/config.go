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
	if _, err := os.Stat(os.ExpandEnv(filepath.Join("$HOME", ".fileCompare"))); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(os.ExpandEnv(filepath.Join("$HOME", ".fileCompare")), 0777)
	}
}

func InitConfigFile() {
	if _, err := os.Stat(os.ExpandEnv(filepath.Join("$HOME", ".fileCompare", ".fileCompare.json"))); errors.Is(err, os.ErrNotExist) {
		os.Create(os.ExpandEnv(filepath.Join("$HOME", ".fileCompare", ".fileCompare.json")))
		initJsonData()
	}
}

func initJsonData() {
	result := FileCompareConfig{}
	stringSlice := make([]string, 1)
	stringSlice[0] = "put data here"

	byteValue, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	// Write back to file
	err = os.WriteFile(os.ExpandEnv(filepath.Join("$HOME", ".fileCompare", ".fileCompare.json")), byteValue, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
