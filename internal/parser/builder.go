package parser

import (
	"fmt"
	"os"
)

type ParserBuilder struct {
	parser Parser
	err    error
}

func NewBuilder(pt string) *ParserBuilder {
	var err error
	if pt != "psql" {
		err = fmt.Errorf("only psql type is supported")
	}

	return &ParserBuilder{
		parser: Parser{
			currentChar:    0x0,
			currentComment: 0x0,
			currentQuote:   0x0,
			inTransaction:  false,
		},
		err: err,
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
	return pb.parser, pb.err
}
