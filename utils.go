package main

import (
	"strings"
	"unicode/utf16"

	"github.com/mattn/go-unicodeclass"
	"github.com/sourcegraph/go-lsp"
)

// WordAt returns word at certain postition
// Stolen from https://github.com/mattn/efm-langserver/blob/master/langserver/handler.go
func (s *server) WordAt(uri lsp.DocumentURI, pos lsp.Position) string {
	text, ok := s.documents[uri]
	if !ok {
		return ""
	}
	lines := strings.Split(text, "\n")
	if pos.Line < 0 || pos.Line > len(lines) {
		return ""
	}
	chars := utf16.Encode([]rune(lines[pos.Line]))
	if pos.Character < 0 || pos.Character > len(chars) {
		return ""
	}
	prevPos := 0
	currPos := -1
	prevCls := unicodeclass.Invalid
	for i, char := range chars {
		currCls := unicodeclass.Is(rune(char))
		if currCls != prevCls {
			if i <= pos.Character {
				prevPos = i
			} else {
				if char == '_' {
					continue
				}
				currPos = i
				break
			}
		}
		prevCls = currCls
	}
	if currPos == -1 {
		currPos = len(chars)
	}
	return string(utf16.Decode(chars[prevPos:currPos]))
}
