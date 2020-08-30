package main

import (
	"regexp"
	"strings"

	"github.com/sourcegraph/go-lsp"
)

var nonEmptyString = regexp.MustCompile(`\S+`)

// WordAt returns word at certain postition
func (s *server) WordAt(uri lsp.DocumentURI, pos lsp.Position) string {
	text, ok := s.documents[uri]
	if !ok {
		return ""
	}
	lines := strings.Split(text, "\n")
	if pos.Line < 0 || pos.Line > len(lines) {
		return ""
	}

	curLine := lines[pos.Line]
	wordIdxs := nonEmptyString.FindAllStringIndex(curLine, -1)
	for _, wordIdx := range wordIdxs {
		if wordIdx[0] <= pos.Character && pos.Character < wordIdx[1] {
			return curLine[wordIdx[0]:wordIdx[1]]
		}
	}

	return ""
}
