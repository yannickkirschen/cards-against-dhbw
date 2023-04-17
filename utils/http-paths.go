package utils

import (
	"errors"
	"strings"
)

func PathParameterFilter(path string, prefix string) (string, error) {
	suffix := strings.TrimPrefix(path, prefix)
	var parameter string
	for _, char := range suffix {
		if char == '/' {
			return "", errors.New("too many parameters (found something after /)")
		}

		parameter += string(char)
	}

	if len(parameter) > 0 {
		return parameter, nil
	}

	return "", errors.New("no parameter found")
}
