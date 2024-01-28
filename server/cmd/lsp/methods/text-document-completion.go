package methods

import (
	"encoding/json"
	"regexp"

	"github.com/baswilson/templ-lsp/cmd/lsp/constants"
	"github.com/baswilson/templ-lsp/cmd/lsp/globals"
)

type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

type CompletionItem struct {
	Label            string             `json:"label"`
	LabelDetails     LabelDetails       `json:"labelDetails"`
	Kind             CompletionItemKind `json:"kind"`
	Detail           string             `json:"detail"`
	Documentation    string             `json:"documentation"`
	InsertText       string             `json:"insertText"`
	InsertTextFormat int                `json:"insertTextFormat"` // 1 = plain text, 2 = snippet
	Tags             []int              `json:"tags"`             // 1 = deprecated
	// this is incomplete
}

type CompletionItemKind int64

const (
	Text        CompletionItemKind = 1
	Method      CompletionItemKind = 2
	Function    CompletionItemKind = 3
	Constructor CompletionItemKind = 4
	Field       CompletionItemKind = 5
	Variable    CompletionItemKind = 6
	Class       CompletionItemKind = 7
	Interface   CompletionItemKind = 8
	Module      CompletionItemKind = 9
	Property    CompletionItemKind = 10
)

type LabelDetails struct {
	Detail      string `json:"detail"`
	Description string `json:"description"`
}

type CompletionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type CompletionMessage struct {
	Message
	Params CompletionParams `json:"params"`
}

func SplitLines(content string) []string {
	lines := []string{}
	currentLine := ""
	for _, char := range content {
		if char == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
			continue
		}
		currentLine += string(char)
	}

	if len(lines) == 0 {
		lines = append(lines, content)
	}

	return lines
}

func isRealHtmlTag(tag string) bool {
	for _, htmlTag := range constants.HtmlTags {
		if tag == htmlTag {
			return true
		}
	}
	return false
}

func attributesByHtmlTag(tag string, attributesToSkip []string) []string {
	allAttributes := []string{}

	for htmlTag, attributes := range constants.HtmlTagsWithAttributes {
		if tag == htmlTag {
			allAttributes = append(attributes, constants.HtmlTagsWithAttributes["*"]...)
		}
	}

	a := []string{}
	for _, attribute := range allAttributes {
		shouldSkip := false
		for _, attributeToSkip := range attributesToSkip {
			if attribute == attributeToSkip {
				shouldSkip = true
				break
			}
		}
		if !shouldSkip {
			a = append(a, attribute)
		}
	}

	return a
}

func filterHtmxAttributes(attributesToSkip []string) map[string]constants.AttributeInfo {
	attributesToReturn := map[string]constants.AttributeInfo{}
	for _, attribute := range constants.HtmxAttributes {
		shouldSkip := false
		for _, attributeToSkip := range attributesToSkip {
			if attribute.Key == attributeToSkip {
				shouldSkip = true
				break
			}
		}
		if !shouldSkip {
			attributesToReturn[attribute.Key] = attribute
		}
	}

	return attributesToReturn
}

func TextDocumentCompletion(rawMessage []byte) (CompletionList, error) {
	var message CompletionMessage
	if err := json.Unmarshal(rawMessage, &message); err != nil {
		WriteLog(err.Error())
		return CompletionList{}, err
	}

	params := message.Params
	content := globals.Documents[params.TextDocument.Uri]
	line := SplitLines(content)[params.Position.Line]
	isInHtmlTag := false
	htmlTag := ""
	indexOfLastOpeningBracket := 0
	indexOfLastClosingBracket := 0

	// check if position is in html tag
	WriteLog("line: " + line)
	for i := 0; i < params.Position.Character; i++ {
		if line[i] == '<' {
			isInHtmlTag = true
			indexOfLastOpeningBracket = i
		} else if line[i] == '>' {
			isInHtmlTag = false
		}
	}

	// find closing bracket position
	for i := params.Position.Character; i < len(line); i++ {
		if line[i] == '>' {
			indexOfLastClosingBracket = i
			break
		}
	}

	regex := regexp.MustCompile(`<([a-zA-Z0-9]+)`)
	if regex.MatchString(line) {
		matches := regex.FindStringSubmatch(line)
		htmlTag = matches[1]
	}

	items := []CompletionItem{}

	if isInHtmlTag {
		isReal := isRealHtmlTag(htmlTag)

		if htmlTag == "" || !isReal {
			for _, tag := range constants.HtmlTags {
				items = append(items, CompletionItem{
					Label:      tag,
					Kind:       Property,
					InsertText: tag + ` `,
				})
			}
		} else if isReal {

			attributesToSkip := []string{}

			if indexOfLastClosingBracket > indexOfLastOpeningBracket {
				regex := regexp.MustCompile(`([a-zA-Z0-9-]+)=`)
				if regex.MatchString(line) {
					matches := regex.FindAllStringSubmatch(line, -1)
					for _, match := range matches {
						attributesToSkip = append(attributesToSkip, match[1])
					}
				}
			}

			for _, attribute := range attributesByHtmlTag(htmlTag, attributesToSkip) {
				items = append(items, CompletionItem{
					Label:            attribute,
					Kind:             Property,
					InsertText:       attribute + `="$1"`,
					InsertTextFormat: 2,
				})
			}

			// add htmx attr
			for _, attribute := range filterHtmxAttributes(attributesToSkip) {
				items = append(items, CompletionItem{
					Label:            attribute.Key,
					Kind:             Property,
					InsertText:       attribute.Key + `="$1"`,
					Detail:           "(HTMX) " + attribute.Label,
					Documentation:    attribute.Description,
					InsertTextFormat: 2,
				})
			}
		}
	}

	return CompletionList{
		IsIncomplete: false,
		Items:        items,
	}, nil
}
