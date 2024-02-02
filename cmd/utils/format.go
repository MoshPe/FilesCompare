package utils

import (
	"fmt"
	"regexp"
	"strings"
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
