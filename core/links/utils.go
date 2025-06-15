package links

import (
	"regexp"
)

// Regex pattern for UTM parameters
var utmParamRegex = regexp.MustCompile(`^utm_`)

// isUTMParam checks if a parameter name matches the UTM pattern
func isUTMParam(param string) bool {
	return utmParamRegex.MatchString(param)
}
