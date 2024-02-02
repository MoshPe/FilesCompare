package config

import (
	"FilesCompare/cmd"
	"FilesCompare/cmd/utils"
	"errors"
	"fmt"
	"github.com/go-test/deep"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

var csvOutput utils.CsvCompareWriter
var referenceFilePath string
var referenceFileName string
var outputPath string

func compareConfigCmd() *cobra.Command {
	var compareCmd = &cobra.Command{
		Use:     "compare-config",
		Short:   "Compare files using the config file",
		Example: `file-compare compare-config`,
		Run: func(cmd *cobra.Command, args []string) {
			fileCompareConfig := viper.AllSettings()
			filesContent := make(map[string]interface{})
			referenceFileName = filepath.Base(referenceFilePath)
			var referenceFileContent interface{}
			fileExt := filepath.Ext(referenceFilePath)
			log.Println("Comparing files of type " + fileExt)

			switch strings.ToLower(fileExt) {
			case utils.Json:
				for _, path := range args {
					filesContent[filepath.Base(path)] = utils.UnmarshalJson(path)
				}
				referenceFileContent = utils.UnmarshalJson(referenceFilePath)
				break
			case utils.YAML, utils.YML:
				for _, path := range args {
					filesContent[filepath.Base(path)] = utils.UnmarshalYaml(path)
				}
				referenceFileContent = utils.UnmarshalYaml(referenceFilePath)
				break
			case utils.XML:
				for _, path := range args {
					filesContent[filepath.Base(path)] = utils.UnmarshalXML(path)
				}
				referenceFileContent = utils.UnmarshalXML(referenceFilePath)
				break
			default:
				panic(errors.New(fmt.Sprintf("invalid file type: %s", fileExt)))
			}

			closeCsv := csvOutput.InitCsv(outputPath, referenceFileName, args)
			defer closeCsv()
			compareFilesContent(referenceFileContent, filesContent)
			compareFilesDates(args)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if referenceFilePath == "" {
				return errors.New("reference file is required")
			}
			return nil
		},
	}
	compareFlags(compareCmd)
	return compareCmd
}

func init() {
	cmd.RootCmd.AddCommand(compareConfigCmd())
}
func compareFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&referenceFilePath, "reference", "r", "", "Reference file to compare")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "compare_results", "Output file name. default: compare_results.csv")
}

func compareFilesContent(reference interface{}, files map[string]interface{}) {
	var diff []string
	for fileName, fileContent := range files {
		diff = deep.Equal(reference, fileContent)
		if diff != nil {
			log.Printf("Files %s to %s are NOT equals", referenceFileName, fileName)
			for _, s := range diff {
				keys, fields := utils.ExtractKeys(s)
				referenceValue, err := getValue(fields, reference)
				if err != nil {
					log.Fatalf("Error extractig field %s from %s", keys, referenceFileName)
				}

				compareToValue, err := getValue(fields, fileContent)
				if err != nil {
					log.Fatalf("Error extractig field %s from %s", keys, fileName)
				}

				csvOutput.WriteRow(referenceValue, compareToValue, keys, fileName)
			}
		} else {
			log.Printf("Files %s to %s are equals", referenceFileName, fileName)
		}
	}
}

func compareFilesDates(files []string) {
	creationReferenceTime, modificationReferenceTime := getCreationModificationTime(referenceFilePath)

	for _, file := range files {
		creationDate, modificationDate := getCreationModificationTime(file)
		csvOutput.WriteRow(creationReferenceTime, creationDate, "Creation Date", file)
		csvOutput.WriteRow(modificationReferenceTime, modificationDate, "Modification Date", file)
	}
}

func getCreationModificationTime(path string) (time.Time, time.Time) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("Error:", err)
		return time.Time{}, time.Time{}
	}
	creationReferenceTime := fileInfo.ModTime()
	modificationReferenceTime := fileInfo.ModTime()
	return creationReferenceTime, modificationReferenceTime
}

func getValue(fields []string, obj interface{}) (interface{}, error) {
	objMap := make(map[string]interface{})
	var ok bool
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Map {
		// Iterate over map keys and values
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			objMap[utils.ConvertToString(key.Interface())] = val.Interface()
		}
	}

	for i, field := range fields {
		val, found := objMap[field]
		if !found {
			return nil, fmt.Errorf("field %s not found", field)
		}

		if i == len(fields)-1 {
			return val, nil
		}

		objMap, ok = val.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid structure")
		}
	}

	return nil, errors.New("invalid path")
}
