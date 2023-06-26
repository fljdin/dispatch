package parser

import (
	"fmt"
	"os"
)

var ParserTypes = []string{"", "sh", "psql"}

type ParserBuilder struct {
	parser Parser
	err    error
}

func NewBuilder(pt string) *ParserBuilder {
	var err error
	var parser Parser

	switch pt {
	case "", "sh":
		parser = &ShParser{}
	case "psql":
		parser = &PsqlParser{
			currentChar:    0x0,
			currentComment: 0x0,
			currentQuote:   0x0,
			inTransaction:  false,
		}
	default:
		err = fmt.Errorf("%s is not supported", pt)
	}

	return &ParserBuilder{
		parser: parser,
		err:    err,
	}
}

func (pb *ParserBuilder) WithContent(content string) *ParserBuilder {
	if pb.err != nil {
		return pb
	}

	pb.parser.SetContent(content)
	return pb
}

func (pb *ParserBuilder) FromFile(filename string) *ParserBuilder {
	data, err := os.ReadFile(filename)
	if err != nil {
		pb.err = fmt.Errorf("error reading file: %w", err)
	}

	pb.parser.SetContent(string(data))
	return pb
}

func (pb *ParserBuilder) Build() (Parser, error) {
	return pb.parser, pb.err
}
