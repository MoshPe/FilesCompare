package utils

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func ExtractKeys(input string) (string, []string) {
	re := regexp.MustCompile(`\[(\w+)\]`)
	matches := re.FindAllStringSubmatch(input, -1)

	var fields []string
	for _, match := range matches {
		fields = append(fields, match[1])
	}
	fullFieldPath := ""
	for _, match := range fields {
		fullFieldPath = fmt.Sprintf("%s.%s", fullFieldPath, match)
	}
	return strings.TrimPrefix(fullFieldPath, "."), fields

}

func ConvertToString(obj any) string {
	return fmt.Sprintf("%v", obj)
}

func GetValue(fields []string, obj interface{}) (interface{}, error) {
	objMap := make(map[string]interface{})
	var ok bool
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Map {
		// Iterate over map keys and values
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			objMap[ConvertToString(key.Interface())] = val.Interface()
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

func GetCreationModificationTime(path string) (time.Time, time.Time) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("Error:", err)
		return time.Time{}, time.Time{}
	}

	//stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	//if !ok {
	//	log.Fatalln("Error: Unable to get file creation time for file", path)
	//}
	//
	//creationReferenceTime := time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
	modificationReferenceTime := fileInfo.ModTime()
	creationReferenceTime := fileInfo.ModTime()
	return creationReferenceTime, modificationReferenceTime
}
