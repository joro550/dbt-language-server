package main

import (
	"regexp"
	"slices"
	"strings"
)

type JinjaParser struct {
	expressionPattern      *regexp.Regexp
	statementPattern       *regexp.Regexp
	commentPattern         *regexp.Regexp
	refPattern             *regexp.Regexp
	macroPattern           *regexp.Regexp
	effectiveJinjaPattern  *regexp.Regexp
	macroDefinitionPattern *regexp.Regexp
}

func NewJinjaParser() JinjaParser {
	expressionPattern := regexp.MustCompile(`{{[\s\S]*?}}`)
	statementPattern := regexp.MustCompile(`{%[\s\S]*?%}`)
	commentPattern := regexp.MustCompile(`{#[\s\S]*?#}`)
	effectiveJinjaPattern := regexp.MustCompile(`{{[\s\S]*?}}|{%[\s\S]*?%}`)
	refPattern := regexp.MustCompile(`{{\s*ref\s*\(\s*['|"](?<project>[a-z_]*?)\s*['|"]\s*(,?\s*['|"](?<model>[a-z_]*?)\s*['|"])?\)\s*}}`)
	macroPattern := regexp.MustCompile(`{{\s*(?<function_name>[a-zA-Z_]*)\s*\([\sA-Za-z=,'_0-9"]*\)\s*}}`)
	macroDefition := regexp.MustCompile(`{%-?\s*macro\s*(?<function_name>[a-zA-Z_]*)\s*\([\sA-Za-z=,'_0-9"]*\)\s*-?%}`)

	return JinjaParser{
		expressionPattern:      expressionPattern,
		statementPattern:       statementPattern,
		commentPattern:         commentPattern,
		refPattern:             refPattern,
		macroPattern:           macroPattern,
		effectiveJinjaPattern:  effectiveJinjaPattern,
		macroDefinitionPattern: macroDefition,
	}
}

func (jp JinjaParser) HasJinjaBlocks(content string) bool {
	return strings.Contains(content, "{") && strings.Contains(content, "}")
}

func (jp JinjaParser) GetJinjaPositions(content string) []Range {
	resultIndicies := jp.effectiveJinjaPattern.FindAllStringIndex(content, -1)

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
	resultIndicies := jp.refPattern.FindAllStringIndex(content, -1)

	if resultIndicies == nil {
		return []ModelReference{}
	}

	modelIndex := jp.refPattern.SubexpIndex("model")
	projectIndex := jp.refPattern.SubexpIndex("project")
	matches := jp.refPattern.FindAllStringSubmatch(content, -1)

	references := []ModelReference{}

	for i, index := range resultIndicies {

		modelName := ""
		if matches[i][modelIndex] != "" {
			modelName = matches[i][modelIndex]
		} else {
			modelName = matches[i][projectIndex]
		}

		references = append(references, ModelReference{
			ModelName: modelName,
			Range:     Range{Start: index[0], End: index[1]},
		})
	}

	return references
}

func (jp JinjaParser) GetMacros(content string) []MacroReference {
	keywords := []string{"ref", "config"}
	resultIndicies := jp.macroPattern.FindAllStringIndex(content, -1)

	if resultIndicies == nil {
		return []MacroReference{}
	}

	macroNames := []MacroReference{}
	functionIndex := jp.macroPattern.SubexpIndex("function_name")
	matches := jp.macroPattern.FindAllStringSubmatch(content, -1)

	for i, match := range matches {

		functionName := match[functionIndex]
		if functionName == "" || slices.Contains(keywords, functionName) {
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

func (jp JinjaParser) GetMacroDefinitions(content string) []MacroReference {
	keywords := []string{"ref", "config"}
	resultIndicies := jp.macroDefinitionPattern.FindAllStringIndex(content, -1)

	if resultIndicies == nil {
		return []MacroReference{}
	}

	macroNames := []MacroReference{}
	functionIndex := jp.macroDefinitionPattern.SubexpIndex("function_name")
	matches := jp.macroDefinitionPattern.FindAllStringSubmatch(content, -1)

	for i, match := range matches {

		functionName := match[functionIndex]
		if functionName == "" || slices.Contains(keywords, functionName) {
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
