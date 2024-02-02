package utils

import (
	"encoding/json"
	"github.com/clbanning/mxj/v2"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const (
	Json string = ".json"
	XML  string = ".xml"
	YAML string = ".yaml"
	YML  string = ".yml"
)

func UnmarshalJson(filepath string) interface{} {
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error in reading file - reading %s - %s", filepath, err)
	}
	var obj interface{}
	err = json.Unmarshal(fileBytes, &obj)
	if err != nil {
		log.Fatalf("Error in unmarshling %s file - reading %s - %s", Json, filepath, err)
	}

	return obj
}

func UnmarshalXML(filepath string) interface{} {
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error in reading file - reading %s - %s", filepath, err)
	}
	obj, err := mxj.NewMapXml(fileBytes)
	if err != nil {
		log.Fatalf("Error in unmarshling %s file - reading %s - %s", XML, filepath, err)
	}

	return obj
}

func UnmarshalYaml(filepath string) interface{} {
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error in reading file - reading %s - %s", filepath, err)
	}
	var obj interface{}
	err = yaml.Unmarshal(fileBytes, &obj)
	if err != nil {
		log.Fatalf("Error in unmarshling %s file - reading %s - %s", YAML, filepath, err)
	}

	return obj
}
