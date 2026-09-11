package regexhandlers

import (
	"net/url"
	"regexp"
	"strings"
)

type Params []string

type QueryPathPattern struct {
	pattern  string
	compiled *regexp.Regexp
	params   Params
}

func NewQueryPathPattern(pattern string, params Params) *QueryPathPattern {
	qpp := new(QueryPathPattern)
	qpp.pattern = pattern
	qpp.mustCompile()
	qpp.params = params
	return qpp
}

func (qpp *QueryPathPattern) normalize() string {
	if !strings.HasSuffix(qpp.pattern, "$") {
		qpp.pattern += `$`
	}
	qpp.pattern = strings.ReplaceAll(qpp.pattern, `{numeric}`, `(\d*)`)
	qpp.pattern = strings.ReplaceAll(qpp.pattern, `{string}`, `([\p{L}|\p{N}|\.|_|\-| ]*)`)
	qpp.pattern = strings.ReplaceAll(qpp.pattern, `{filename}`, `([\p{L}|\p{N}|\.|_|\-| ]*)`)
	qpp.pattern = strings.ReplaceAll(qpp.pattern, `{any}`, `([^/]*)`)
	return qpp.pattern
}

func (qpp *QueryPathPattern) mustCompile() {
	qpp.compiled = regexp.MustCompile(qpp.normalize())
}

func (qpp *QueryPathPattern) match(url string) bool {
	return qpp.compiled.MatchString(url)
}

func (qpp *QueryPathPattern) GetValues(urlPath string) url.Values {
	parsedParamsInURL := qpp.compiled.FindAllStringSubmatch(urlPath, -1)
	paramsValues := make(url.Values)
	for _, submatch := range parsedParamsInURL {
		for orderIndex, value := range submatch {
			if orderIndex == 0 {
				continue
			}
			paramsValues[qpp.params[orderIndex-1]] = append(paramsValues[qpp.params[orderIndex-1]], value)
		}
	}
	return paramsValues
}
