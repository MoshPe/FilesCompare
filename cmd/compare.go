package cmd

import (
	"FilesCompare/cmd/utils"
	"errors"
	"fmt"
	"github.com/go-test/deep"
	"github.com/spf13/cobra"
	_ "gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

var ReferenceFilePath string
var ReferenceFileName string
var OutputPath string
var csvOutput utils.CsvCompareWriter

func compareCmd() *cobra.Command {
	var compareCmd = &cobra.Command{
		Use:     "compare [...PATHS]",
		Short:   "Compare files to a reference file",
		Example: `file-compare compare --reference /path/to/reference/file --files /path/1 /path/2 /path/3`,
		Run: func(cmd *cobra.Command, args []string) {

			filesContent := make(map[string]interface{})
			ReferenceFileName = filepath.Base(ReferenceFilePath)
			var referenceFileContent interface{}
			fileExt := filepath.Ext(ReferenceFilePath)
			log.Println("Comparing files of type " + fileExt)

			switch strings.ToLower(fileExt) {
			case utils.Json:
				for _, path := range args {
					filesContent[filepath.Base(path)] = utils.UnmarshalJson(path)
				}
				referenceFileContent = utils.UnmarshalJson(ReferenceFilePath)
				break
			case utils.YAML, utils.YML:
				for _, path := range args {
					filesContent[filepath.Base(path)] = utils.UnmarshalYaml(path)
				}
				referenceFileContent = utils.UnmarshalYaml(ReferenceFilePath)
				break
			case utils.XML:
				for _, path := range args {
					filesContent[filepath.Base(path)] = utils.UnmarshalXML(path)
				}
				referenceFileContent = utils.UnmarshalXML(ReferenceFilePath)
				break
			default:
				panic(errors.New(fmt.Sprintf("invalid file type: %s", fileExt)))
			}

			closeCsv := csvOutput.InitCsv(OutputPath, ReferenceFileName, args)
			defer closeCsv()
			compareFilesContent(referenceFileContent, filesContent)
			compareFilesDates(args)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("required at least 1 file to compare")
			}
			if ReferenceFilePath == "" {
				return errors.New("reference file is required")
			}
			return nil
		},
	}
	compareFlags(compareCmd)
	return compareCmd
}

func init() {
	RootCmd.AddCommand(compareCmd())
}
func compareFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&ReferenceFilePath, "reference", "r", "", "Reference file to compare")
	cmd.Flags().StringVarP(&OutputPath, "output", "o", "compare_results", "Output file name")
}

func compareFilesContent(reference interface{}, files map[string]interface{}) {
	var diff []string
	for fileName, fileContent := range files {
		diff = deep.Equal(reference, fileContent)
		if diff != nil {
			log.Printf("Files %s to %s are NOT equals", ReferenceFileName, fileName)
			for _, s := range diff {
				keys, fields := utils.ExtractKeys(s)
				referenceValue, err := getValue(fields, reference)
				if err != nil {
					log.Fatalf("Error extractig field %s from %s", keys, ReferenceFileName)
				}

				compareToValue, err := getValue(fields, fileContent)
				if err != nil {
					log.Fatalf("Error extractig field %s from %s", keys, fileName)
				}

				csvOutput.WriteRow(referenceValue, compareToValue, keys, fileName)
			}
		} else {
			log.Printf("Files %s to %s are equals", ReferenceFileName, fileName)
		}
	}
}

func compareFilesDates(files []string) {
	creationReferenceTime, modificationReferenceTime := getCreationModificationTime(ReferenceFilePath)

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
