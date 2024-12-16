package util

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

func ValidPath(path string) error {
	re := regexp.MustCompile(`\.\.`)
	if re.MatchString(path) {
		return errors.New("can't have '..' in URL path")
	}
	return nil
}

func CleanPath(path string) (string, error) {
	var err error

	err = ValidPath(path)
	if err != nil {
		return "", err
	}

	if path != "/" {
		cleanPath, found := strings.CutPrefix(path, "/")
		if !found {
			return "", errors.New("expected URL path to start with '/'")
		}
		return filepath.Clean(cleanPath), nil
	}

	return "index.html", nil
}
