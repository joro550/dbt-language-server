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
	refPattern := regexp.MustCompile(`{{\s*ref\s*\(\s*['|"](?<project>[a-z_]*?)\s*['|"]\s*(,?\s*['|"](?<model>[a-z_]*?)\s*['|"])?\)\s*}}`)

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
	resultIndicies := jp.refPattern.FindAllIndex(byteContent, -1)

	if resultIndicies == nil {
		return []ModelReference{}
	}

	nameGroups := []map[string]string{}
	names := jp.refPattern.SubexpNames()
	matches := jp.refPattern.FindAllStringSubmatch(string(byteContent), -1)

	for i := range matches {
		matchMap := make(map[string]string)

		for j, name := range names {
			if j != 0 && name != "" {
				matchMap[name] = matches[i][j]
			}
		}

		nameGroups = append(nameGroups, matchMap)
	}

	references := []ModelReference{}
	for i, index := range resultIndicies {

		modelName := ""
		if nameGroups[i]["model"] != "" {
			modelName = nameGroups[i]["model"]
		} else {
			modelName = nameGroups[i]["project"]
		}

		references = append(references, ModelReference{
			ModelName: modelName,
			Range:     Range{Start: index[0], End: index[1]},
		})
	}

	return references
}
