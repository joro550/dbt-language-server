package main

import (
	"regexp"
	"slices"
	"strings"
)

type JinjaParser struct {
	expressionPattern     *regexp.Regexp
	statementPattern      *regexp.Regexp
	commentPattern        *regexp.Regexp
	refPattern            *regexp.Regexp
	macroPattern          *regexp.Regexp
	effectiveJinjaPattern *regexp.Regexp
}

func NewJinjaParser() JinjaParser {
	expressionPattern := regexp.MustCompile(`{{[\s\S]*?}}`)
	statementPattern := regexp.MustCompile(`{%[\s\S]*?%}`)
	commentPattern := regexp.MustCompile(`{#[\s\S]*?#}`)
	effectiveJinjaPattern := regexp.MustCompile(`{{[\s\S]*?}}|{%[\s\S]*?%}`)
	refPattern := regexp.MustCompile(`{{\s*ref\s*\(\s*['|"](?<project>[a-z_]*?)\s*['|"]\s*(,?\s*['|"](?<model>[a-z_]*?)\s*['|"])?\)\s*}}`)
	macroPattern := regexp.MustCompile(`{{\s*(?<function_name>[a-zA-Z_]*)\s*\([\s\S]*\)\s*}}`)

	return JinjaParser{
		expressionPattern:     expressionPattern,
		statementPattern:      statementPattern,
		commentPattern:        commentPattern,
		refPattern:            refPattern,
		macroPattern:          macroPattern,
		effectiveJinjaPattern: effectiveJinjaPattern,
	}
}

func (jp JinjaParser) HasJinjaBlocks(content string) bool {
	return strings.Contains(content, "{") && strings.Contains(content, "}")
}

func (jp JinjaParser) GetJinjaPositions(content string) []Range {
	byteContent := []byte(content)
	resultIndicies := jp.refPattern.FindAllIndex(byteContent, -1)

	if resultIndicies == nil {
		return []Range{}
	}

	ranges := []Range{}
	for _, result := range resultIndicies {
		ranges = append(ranges, Range{Start: result[0], End: result[1]})
	}
	return ranges
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

func (jp JinjaParser) GetMacros(content string) []MacroReference {
	keywords := []string{"ref"}
	byteContent := []byte(content)
	resultIndicies := jp.macroPattern.FindAllIndex(byteContent, -1)

	if resultIndicies == nil {
		return []MacroReference{}
	}

	macroNames := []MacroReference{}
	names := jp.macroPattern.SubexpNames()
	matches := jp.macroPattern.FindAllStringSubmatch(string(byteContent), -1)

	for i := range matches {
		matchMap := make(map[string]string)

		for j, name := range names {
			if j != 0 && name != "" {
				matchMap[name] = matches[i][j]
			}
		}

		functionName, ok := matchMap["function_name"]
		if !ok || functionName == "" || slices.Contains(keywords, functionName) {
			continue
		}

		macroNames = append(macroNames, MacroReference{
			ModelName: functionName,
			Range: Range{
				Start: resultIndicies[i][0],
				End:   resultIndicies[i][1],
			},
		})

	}

	return macroNames
}
