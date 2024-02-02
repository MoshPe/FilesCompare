package pkg

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

type FileCompareConfig struct {
	Ips                []string `json:"ips"`
	InstallationFolder string   `json:"installation_folder"`
	ExcludePatterns    []string `json:"exclude_patterns"`
}

func InitConfigDir() {
	if _, err := os.Stat(os.ExpandEnv(filepath.Join("$HOME", ".fileCompare"))); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(os.ExpandEnv(filepath.Join("$HOME", ".fileCompare")), 0755)
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
	result.Ips = stringSlice
	result.ExcludePatterns = stringSlice

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
