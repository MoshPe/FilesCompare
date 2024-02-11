package config

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func FindFilesMatchingPattern(dir, pattern string) ([]string, error) {
	matchingFiles := make([]string, 0)

	// Compile pattern into regexp
	regexpPattern := regexp.QuoteMeta(pattern)
	regexpPattern = "^" + regexpPattern + "$"
	regexpPattern = strings.ReplaceAll(regexpPattern, `\*`, `.*`)
	regexpPattern = strings.ReplaceAll(regexpPattern, `\?`, `.`)

	re, err := regexp.Compile(regexpPattern)
	if err != nil {
		return nil, err
	}

	// Walk through the directory and its subdirectories
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the file matches the pattern
		if re.MatchString(info.Name()) {
			matchingFiles = append(matchingFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return matchingFiles, nil
}

func FindFilesMatchingPatterns(dir string, patterns []string) ([]string, error) {
	matchingFiles := make([]string, 0)

	// Compile patterns into regexps
	regexps := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		regexpPattern := regexp.QuoteMeta(pattern)
		regexpPattern = "^" + regexpPattern + "$"
		regexpPattern = strings.ReplaceAll(regexpPattern, `\*`, `.*`)
		regexpPattern = strings.ReplaceAll(regexpPattern, `\?`, `.`)
		re, err := regexp.Compile(regexpPattern)
		if err != nil {
			return nil, err
		}
		regexps[i] = re
	}

	// Walk through the directory and its subdirectories
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the file matches any of the patterns
		for _, re := range regexps {
			if re.MatchString(info.Name()) {
				matchingFiles = append(matchingFiles, path)
				break // Move to the next file after finding a match
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return matchingFiles, nil
}
