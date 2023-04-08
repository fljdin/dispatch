package parser

import (
	"fmt"
	"os"
)

type ParserBuilder struct {
	parser Parser
	err    error
}

func NewParserBuilder(pt string) *ParserBuilder {
	return &ParserBuilder{
		parser: Parser{
			Type:           pt,
			currentChar:    0x0,
			currentComment: 0x0,
			currentQuote:   0x0,
			inTransaction:  false,
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
