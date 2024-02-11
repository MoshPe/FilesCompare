package config

import (
	"FilesCompare/cmd"
	"FilesCompare/cmd/utils"
	"FilesCompare/pkg"
	"github.com/go-test/deep"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "gopkg.in/yaml.v3"
	"log"
	"path/filepath"
)

type Application struct {
	FilesContent map[string]FileContent
}

type FileContent struct {
	Content interface{}
	path    string
}

var csvOutput utils.CsvCompareWriter
var referenceApplicationName string
var outputPath string
var FileCompareConfig pkg.FileCompareConfig
var applicationsFiles map[string]Application
var referenceFilesContent map[string]FileContent

func compareConfigCmd() *cobra.Command {
	var compareCmd = &cobra.Command{
		Use:     "compare-config",
		Short:   "Compare files using the config file",
		Example: `file-compare compare-config`,
		Run: func(cmd *cobra.Command, args []string) {
			config := viper.AllSettings()

			if err := mapstructure.Decode(config, &FileCompareConfig); err != nil {
				log.Fatalln("Error decoding settings:", err)
			}

			getFilesPaths()

			keys := make([]string, 0, len(applicationsFiles))
			for key := range applicationsFiles {
				keys = append(keys, key)
			}
			closeCsv := csvOutput.InitCsvWithAppNames(outputPath, FileCompareConfig.FilesToCompare, keys)
			defer closeCsv()
			compareFilesDates()
			prepFilesToCompare()
		},
	}
	compareFlags(compareCmd)
	return compareCmd
}

func init() {
	cmd.RootCmd.AddCommand(compareConfigCmd())
}
func compareFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&outputPath, "output", "o", "compare_results", "Output file name. default: compare_results.csv")
}

func compareFilesContent(reference interface{}, files map[string]interface{}, fileName string) {
	var diff []string
	for appName, fileContent := range files {
		diff = deep.Equal(reference, fileContent)
		if diff != nil {
			log.Printf("File %s are NOT equals - AppName : %s", fileName, appName)
			for _, s := range diff {
				keys, fields := utils.ExtractKeys(s)
				referenceValue, err := utils.GetValue(fields, reference)
				if err != nil {
					log.Fatalf("Error extractig field %s from reference file %s", keys, fileName)
				}

				compareToValue, err := utils.GetValue(fields, fileContent)
				if err != nil {
					log.Fatalf("Error extractig field %s from application file %s - AppName : %s", keys, fileName, appName)
				}

				csvOutput.WriteRow(referenceValue, compareToValue, keys, fileName)
			}
		} else {
			log.Printf("Files %s are equals - AppName : %s", fileName, appName)
		}
	}
}

func getFilesPaths() {
	referenceFilesContent = make(map[string]FileContent)
	referenceApplicationName = filepath.Base(FileCompareConfig.ReferenceApplication)

	referencePatterns, err := FindFilesMatchingPatterns(FileCompareConfig.ReferenceApplication, FileCompareConfig.FilesToCompare)
	if err != nil {
		log.Fatalln("Error iterating directories:", err)
	}
	for _, pattern := range referencePatterns {
		referenceFilesContent[filepath.Base(pattern)] = FileContent{
			utils.MarshalFile(pattern),
			pattern,
		}
	}

	applicationsFiles = make(map[string]Application)
	for _, application := range FileCompareConfig.CompareApplications {
		var applicationFileContent Application
		applicationFileContent.FilesContent = make(map[string]FileContent)
		patterns, err := FindFilesMatchingPatterns(application.Path, FileCompareConfig.FilesToCompare)
		if err != nil {
			log.Fatalln("Error iterating directories:", err)
		}
		for _, pattern := range patterns {
			applicationFileContent.FilesContent[filepath.Base(pattern)] = FileContent{
				utils.MarshalFile(pattern),
				pattern,
			}
		}
		applicationsFiles[application.ApplicationName] = applicationFileContent
	}
	log.Println("Done organizing files")
}

func prepFilesToCompare() {
	for _, pattern := range FileCompareConfig.FilesToCompare {
		filesToCompare := make(map[string]interface{})
		for applicationName, application := range applicationsFiles {
			filesToCompare[applicationName] = application.FilesContent[pattern].Content
		}
		compareFilesContent(referenceFilesContent[pattern].Content, filesToCompare, pattern)
	}
}

func compareFilesDates() {
	for _, application := range applicationsFiles {
		for pattern, fileContent := range application.FilesContent {
			creationReferenceTime, modificationReferenceTime := utils.GetCreationModificationTime(referenceFilesContent[pattern].path)
			creationDate, modificationDate := utils.GetCreationModificationTime(fileContent.path)
			csvOutput.WriteRow(creationReferenceTime, creationDate, "Creation Date", pattern)
			csvOutput.WriteRow(modificationReferenceTime, modificationDate, "Modification Date", pattern)
		}
	}
}
