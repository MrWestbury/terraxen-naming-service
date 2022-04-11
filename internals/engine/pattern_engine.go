package engine

import (
	"regexp"
	"strings"
)

func ResolvePattern(pattern string, attribs map[string]string) (string, error) {
	var re = regexp.MustCompile(`{(.*?)}`)

	patternVars := re.FindAllString(pattern, -1)

	result := pattern
	for _, ptVar := range patternVars {
		matches := re.FindStringSubmatch(ptVar)
		rep, found := attribs[matches[1]]
		if found {
			result = strings.ReplaceAll(result, ptVar, rep)
		}
	}
	return result, nil
}
