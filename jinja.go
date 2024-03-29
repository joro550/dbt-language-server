package main

import (
	"regexp"
	"strings"
)

type JinjaParser struct {
	expressionPattern *regexp.Regexp
	statementPattern  *regexp.Regexp
	commentPattern    *regexp.Regexp
	refPattern        *regexp.Regexp
}

func NewJinjaParser() JinjaParser {
	expressionPattern, _ := regexp.Compile(`{{[\s\S]*?}}`)
	statementPattern, _ := regexp.Compile(`{%[\s\S]*?%}`)
	commentPattern, _ := regexp.Compile(`{#[\s\S]*?#}`)
	refPattern, _ := regexp.Compile(`{{\s*ref\s*\(\s*(?<start_quote>['|"])(.*?)\\k<start_quote>\s*\)\s*}}`)

	return JinjaParser{
		expressionPattern: expressionPattern,
		statementPattern:  statementPattern,
		commentPattern:    commentPattern,
		refPattern:        refPattern,
	}
}

func (jp JinjaParser) HasJinjaBlocks(content string) bool {
	return strings.Contains(content, "{") && strings.Contains(content, "}")
}

func (jp JinjaParser) GetAllRefTags(content string) []ModelReference {
	byteContent := []byte(content)
	resultText := jp.refPattern.FindAll(byteContent, -1)
	resultIndicies := jp.refPattern.FindAllIndex(byteContent, -1)

	if resultIndicies == nil {
		return []ModelReference{}
	}

	references := []ModelReference{}
	for i, index := range resultIndicies {
		modelName := string(resultText[i])
		references = append(references, ModelReference{
			ModelName: modelName,
			Range:     Range{Start: index[0], End: index[0] + index[1]},
		})
	}

	return references
}
