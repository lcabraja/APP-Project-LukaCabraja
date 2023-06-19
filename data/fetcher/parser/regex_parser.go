package parser

import (
	"github.com/lcabraja/APP-Project-LukaCabraja/log"
	"regexp"
)

type RegexParser struct {
	pattern  *regexp.Regexp
	submatch int
}

func NewRegexParser(pattern *regexp.Regexp, submatch int) *RegexParser {
	return &RegexParser{
		pattern:  pattern,
		submatch: submatch,
	}
}

func (rp *RegexParser) Parse(data []byte) (string, error) {
	results := rp.pattern.FindStringSubmatch(string(data))
	log.Dev("Regex Found: %v\n", results)
	return results[0], nil
}
