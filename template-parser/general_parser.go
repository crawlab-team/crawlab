package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
	"regexp"
	"strings"
)

type GeneralParser struct {
	tagPattern   string
	tagRegexp    *regexp.Regexp
	mathPattern  string
	mathRegexp   *regexp.Regexp
	template     string
	indexes      [][]int
	matches      [][]string
	placeholders []string
	variables    []Variable
	vm           *otto.Otto
}

const VariableNameResult = "result"
const ValueNameNA = "N/A"

func (p *GeneralParser) Parse(template string) (err error) {
	p.template = template
	p.indexes = p.tagRegexp.FindAllStringIndex(template, -1)
	p.matches = p.tagRegexp.FindAllStringSubmatch(template, -1)
	for _, arr := range p.matches {
		p.placeholders = append(p.placeholders, arr[1])
	}

	return nil
}

func (p *GeneralParser) Render(args ...interface{}) (content string, err error) {
	// render tag content
	content, err = p.renderTagContent(args...)
	if err != nil {
		return content, err
	}

	// render math content
	content, err = p.renderMathContent(content)
	if err != nil {
		return content, err
	}

	return content, nil
}

func (p *GeneralParser) renderTagContent(args ...interface{}) (content string, err error) {
	// validate
	if len(args) == 0 {
		return "", errors.New("no arguments")
	}

	// first argument
	arg := args[0]

	// content
	content = p.template

	// old strings
	var oldStrList []string
	for _, index := range p.indexes {
		// old string
		oldStr := content[index[0]:index[1]]
		oldStrList = append(oldStrList, oldStr)
	}

	// iterate placeholders
	for i, placeholder := range p.placeholders {
		// variable
		v, err := NewVariable(arg, placeholder)
		if err != nil {
			return "", err
		}

		// value
		value, err := v.GetValue()
		if err != nil || value == nil {
			value = ValueNameNA
		}

		// old string
		oldStr := oldStrList[i]

		// new string
		var newStr string
		switch value.(type) {
		case string:
			newStr = value.(string)
		default:
			newStrBytes, err := json.Marshal(value)
			if err != nil {
				return "", err
			}
			newStr = string(newStrBytes)
		}

		// replace old string with new string
		content = strings.Replace(content, oldStr, newStr, 1)
	}

	return content, nil
}

func (p *GeneralParser) renderMathContent(inputContent string) (content string, err error) {
	content = inputContent
	indexes := p.mathRegexp.FindAllStringIndex(inputContent, -1)
	matches := p.mathRegexp.FindAllStringSubmatch(inputContent, -1)

	// old strings
	var oldStrList []string
	for _, index := range indexes {
		// old string
		oldStr := content[index[0]:index[1]]
		oldStrList = append(oldStrList, oldStr)
	}

	// iterate matches
	for i, m := range matches {
		// js script to run to get evaluate result
		script := fmt.Sprintf("%s = %s; %s", VariableNameResult, m[1], VariableNameResult)

		// replace NA
		script = strings.ReplaceAll(script, ValueNameNA, "NaN")

		// value
		value, err := p.vm.Run(script)
		if err != nil {
			return "", err
		}

		// old string
		oldStr := oldStrList[i]

		// new string
		newStr := value.String()

		// replace old string with new string
		content = strings.Replace(content, oldStr, newStr, 1)
	}

	return content, nil
}

func (p *GeneralParser) GetPlaceholders() (placeholders []string) {
	return p.placeholders
}

func NewGeneralParser() (p Parser, err error) {
	// tag regexp
	tagPrefix := "\\{\\{"
	tagSuffix := "\\}\\}"
	tagBasicChars := "\\$\\.\\w_"
	tagAssociateChars := "\\[\\]:"
	tagPattern := fmt.Sprintf(
		"%s *([%s%s]+) *%s",
		tagPrefix,
		tagBasicChars,
		tagAssociateChars,
		tagSuffix,
	)
	tagRegexp, err := regexp.Compile(tagPattern)
	if err != nil {
		return nil, err
	}

	// math regexp
	mathPrefix := "\\{#"
	mathSuffix := "#\\}"
	mathBasicChars := " \\(\\)"
	mathOpChars := "\\+\\-\\*/%"
	mathNumChars := "\\d\\."
	mathSpecialChars := "(?:NaN|null)"
	mathPattern := fmt.Sprintf(
		"%s *([%s%s%s%s]+) *%s",
		mathPrefix,
		mathBasicChars,
		mathOpChars,
		mathNumChars,
		mathSpecialChars,
		mathSuffix,
	)
	mathRegexp, err := regexp.Compile(mathPattern)
	if err != nil {
		return nil, err
	}

	// math vm
	vm := otto.New()

	// parser
	p = &GeneralParser{
		tagPattern:  tagPattern,
		tagRegexp:   tagRegexp,
		mathPattern: mathPattern,
		mathRegexp:  mathRegexp,
		vm:          vm,
	}

	return p, nil
}
