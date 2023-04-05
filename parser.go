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
	var inLiteral bool

	for i := 0; i < len(p.Content); i++ {
		char := p.Content[i]
		currentQuery.WriteRune(rune(char))

		if currentQuery.Len() > 0 && char == ';' && !inLiteral {
			commands = append(commands, currentQuery.String())
			currentQuery.Reset()
		} else if char == '\'' || char == '"' {
			inLiteral = !inLiteral
		}
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
