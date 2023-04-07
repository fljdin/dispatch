package parser

import (
	"fmt"
	"strings"
)

var ParserTypes = []string{"psql"}

type Parser struct {
	Type    string
	Content string

	currentQuery strings.Builder
}

func (p *Parser) VerifyType() error {
	for _, pt := range ParserTypes {
		if p.Type == pt {
			return nil
		}
	}
	return fmt.Errorf("%s is an invalid type for parsing file", p.Type)
}

func (p *Parser) isBeginPattern() bool {
	queryLen := p.currentQuery.Len()

	if queryLen >= 5 {
		beginWord := p.currentQuery.String()[queryLen-5:]
		return beginWord == "BEGIN"
	}

	return false
}

func (p *Parser) isEndPattern() bool {
	queryLen := p.currentQuery.Len()
	var patternFound bool = false

	if queryLen >= 3 {
		endWord := p.currentQuery.String()[queryLen-3:]
		patternFound = patternFound || (endWord == "END")
	}

	if queryLen >= 6 {
		endWord := p.currentQuery.String()[queryLen-6:]
		patternFound = patternFound || (endWord == "COMMIT")
	}

	if queryLen >= 8 {
		endWord := p.currentQuery.String()[queryLen-8:]
		patternFound = patternFound || (endWord == "ROLLBACK")
	}

	return patternFound
}

func (p *Parser) Parse() []string {
	var commands []string
	var currentChar byte
	var previousChar byte = 0x0
	var currentQuote byte = 0x0
	var currentComment byte = 0x0
	var inTrxBlock bool = false

	for i := 0; i < len(p.Content); i++ {
		currentChar = p.Content[i]
		currentSymbols := string(previousChar) + string(currentChar)
		p.currentQuery.WriteRune(rune(currentChar))

		// return complete query
		if p.currentQuery.Len() > 0 && currentChar == ';' &&
			currentComment == 0x0 && currentQuote == 0x0 && !inTrxBlock {

			commands = append(commands, p.currentQuery.String())
			p.currentQuery.Reset()
			continue
		}

		// ignore semicolon in litterals
		if currentChar == '\'' || currentChar == '"' || currentSymbols == "$$" {
			if currentQuote == 0x0 {
				currentQuote = currentChar
			} else if currentQuote == currentChar {
				currentQuote = 0x0
			}
		}

		// ignore semicolon in comments
		if currentComment == 0x0 &&
			(currentSymbols == "--" || currentSymbols == "/*") {

			currentComment = currentChar // opening a comment
		}

		if (currentComment == '-' && currentChar == '\n') ||
			(currentComment == '*' && currentSymbols == "*/") {

			currentComment = 0x0 // ending a comment
		}

		// ignore semicolon in transaction block
		var isKeyWord bool = (currentComment == 0x0 && currentQuote == 0x0)

		if p.isBeginPattern() && isKeyWord {
			inTrxBlock = true
		}

		if p.isEndPattern() && isKeyWord {
			inTrxBlock = false
		}

		previousChar = currentChar
	}

	return commands
}
