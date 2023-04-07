package parser

import (
	"fmt"
	"strings"
)

var ParserTypes = []string{"psql"}

type Parser struct {
	Type    string
	Content string
}

func (p Parser) VerifyType() error {
	for _, pt := range ParserTypes {
		if p.Type == pt {
			return nil
		}
	}
	return fmt.Errorf("%s is an invalid type for parsing file", p.Type)
}

func (p *Parser) Parse() []string {
	var commands []string
	var currentQuery strings.Builder
	var currentChar byte
	var previousChar byte = 0x0
	var currentQuote byte = 0x0
	var currentComment byte = 0x0
	var inTrxBlock bool = false

	for i := 0; i < len(p.Content); i++ {
		currentChar = p.Content[i]
		currentSymbols := string(previousChar) + string(currentChar)
		currentQuery.WriteRune(rune(currentChar))

		// return complete query
		if currentQuery.Len() > 0 && currentChar == ';' &&
			currentComment == 0x0 && currentQuote == 0x0 && !inTrxBlock {

			commands = append(commands, currentQuery.String())
			currentQuery.Reset()
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
		var beginWord string
		var endWord string
		var rollbackWord string
		var isKeyWord bool = (currentComment == 0x0 && currentQuote == 0x0)

		if i >= 4 {
			beginWord = string(p.Content[i-4:i]) + string(currentChar)

			if beginWord == "BEGIN" && isKeyWord {
				inTrxBlock = true
			}
		}

		if i >= 2 {
			endWord = string(p.Content[i-2:i]) + string(currentChar)

			if endWord == "END" && isKeyWord {
				inTrxBlock = false
			}
		}

		if i >= 5 {
			endWord = string(p.Content[i-5:i]) + string(currentChar)

			if endWord == "COMMIT" && isKeyWord {
				inTrxBlock = false
			}
		}

		if i >= 7 {
			rollbackWord = string(p.Content[i-7:i]) + string(currentChar)

			if rollbackWord == "ROLLBACK" && isKeyWord {
				inTrxBlock = false
			}
		}

		previousChar = currentChar
	}

	return commands
}
