package links

import (
	"regexp"
	"strings"
)

// Regex pattern for UTM parameters
var utmParamRegex = regexp.MustCompile(`^utm_`)

// isUTMParam checks if a parameter name matches the UTM pattern
func isUTMParam(param string) bool {
	return utmParamRegex.MatchString(param)
}

func addTrailingSlash(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	}
	return path + "/"
}

func stripTrailingSlash(path string) string {
	if len(path) <= 1 {
		return path
	}
	return strings.TrimRight(path, "/")
}
