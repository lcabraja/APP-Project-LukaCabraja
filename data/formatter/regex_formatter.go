package formatter

import (
	"fmt"
	"github.com/lcabraja/APP-Project-LukaCabraja/data/map_manager"
	"regexp"
)

type RegexFormatter struct {
	mm  map_manager.MapManager
	vt  map_manager.ValueType
	key string

	regexInput  string
	regexOutput string
	input       *regexp.Regexp
	output      *regexp.Regexp
}

func NewRegexFormatter(regexInput string, regexOutput string) (*RegexFormatter, error) {
	input, err := regexp.Compile(regexInput)
	if err != nil {
		return nil, err
	}

	output, err := regexp.Compile(regexInput)
	if err != nil {
		return nil, err
	}

	return &RegexFormatter{
		regexInput:  regexInput,
		regexOutput: regexOutput,
		input:       input,
		output:      output,
	}, nil
}

func (rf *RegexFormatter) Format(value interface{}) (interface{}, error) {
	strval := fmt.Sprintf("%v", value)
	return rf.FormatRegex(strval), nil
}

func (rf *RegexFormatter) FormatRegex(value string) string {
	return rf.output.ReplaceAllString(value, rf.regexOutput)
}
