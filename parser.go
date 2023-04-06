package main

import (
	"fmt"
	"os"
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

	for i := 0; i < len(p.Content); i++ {
		currentChar = p.Content[i]
		currentSymbols := string(previousChar) + string(currentChar)
		currentQuery.WriteRune(rune(currentChar))

		// return complete query
		if currentQuery.Len() > 0 && currentChar == ';' &&
			currentComment == 0x0 && currentQuote == 0x0 {

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

		previousChar = currentChar
	}

	return commands
}

type ParserBuilder struct {
	parser Parser
	err    error
}

func NewParserBuilder(pt string) *ParserBuilder {
	return &ParserBuilder{
		parser: Parser{
			Type: pt,
		},
	}
}

func (pb *ParserBuilder) WithContent(content string) *ParserBuilder {
	pb.parser.Content = content
	return pb
}

func (pb *ParserBuilder) FromFile(filename string) *ParserBuilder {
	data, err := os.ReadFile(filename)
	if err != nil {
		pb.err = fmt.Errorf("error reading file: %w", err)
	}

	pb.parser.Content = pb.WithContent(string(data)).parser.Content
	return pb
}

func (pb *ParserBuilder) Build() (Parser, error) {
	if err := pb.parser.VerifyType(); err != nil {
		pb.err = err
	}

	return pb.parser, pb.err
}
