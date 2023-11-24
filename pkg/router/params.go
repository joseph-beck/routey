package router

import (
	"errors"
	"strings"
)

func parseParams(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	segments := strings.Split(s, "/")
	pattern := ""

	for _, segment := range segments {
		if segment != "" {
			if pattern != "" {
				pattern += "/"
			}

			if strings.HasPrefix(segment, ":") {
				p := segment[1:]
				pattern += "(?P<" + p + ">\\w+)"
			} else {
				return "", errors.New("bad parameters found in the params")
			}
		}
	}

	pattern = "/" + pattern

	return pattern, nil
}

func parsePathParams(s string) (string, error) {
	if s == "" {
		return "/", nil
	}

	segments := strings.Split(s, "/")
	pattern := ""

	for _, segment := range segments {
		if segment != "" {
			pattern += "/"

			if strings.HasPrefix(segment, ":") {
				p := segment[1:]
				pattern += "(?P<" + p + ">\\w+)"
			} else {
				pattern += segment
			}

			continue
		}
	}

	return pattern, nil
}
