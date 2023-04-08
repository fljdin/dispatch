package parser

import (
	"fmt"
	"strings"
)

var ParserTypes = []string{"psql"}

type Parser struct {
	Type    string
	Content string

	currentQuery   strings.Builder
	currentChar    byte
	currentComment byte // either - or *
	currentQuote   byte // either ', " or $
	inTransaction  bool
}

func (p *Parser) VerifyType() error {
	for _, pt := range ParserTypes {
		if p.Type == pt {
			return nil
		}
	}
	return fmt.Errorf("%s is an invalid type for parsing file", p.Type)
}

func (p *Parser) inComment() bool {
	return p.currentComment != 0x0
}

func (p *Parser) handleComments() {
	// opening a comment
	if !p.inComment() && (p.match("--") || p.match("/*")) {
		p.currentComment = p.currentChar
	}

	// ending a comment
	if (p.currentComment == '-' && p.currentChar == '\n') ||
		(p.currentComment == '*' && p.match("*/")) {
		p.currentComment = 0x0
	}
}

func (p *Parser) inQuotedLiteral() bool {
	return p.currentQuote != 0x0
}

func (p *Parser) handleLiterals() {
	if p.match("'", "\"", "$$") {
		if !p.inQuotedLiteral() {
			// enter into litteral
			p.currentQuote = p.currentChar
		} else if p.currentQuote == p.currentChar {
			// exit a litteral
			p.currentQuote = 0x0
		}
	}
}

func (p *Parser) inCommentOrLiteral() bool {
	return p.inComment() || p.inQuotedLiteral()
}

func (p *Parser) handleTransactions() {
	if p.match("BEGIN") && !p.inCommentOrLiteral() {
		p.inTransaction = true
	}

	if p.match("END", "COMMIT", "ROLLBACK") && !p.inCommentOrLiteral() {
		p.inTransaction = false
	}
}

func (p *Parser) isQueryComplete() bool {
	return p.currentChar == ';' &&
		!p.inCommentOrLiteral() &&
		!p.inTransaction
}

func (p *Parser) match(pattern ...string) bool {
	queryLen := p.currentQuery.Len()
	patternFound := false

	for _, pattern := range pattern {
		if patternLen := len(pattern); queryLen >= patternLen {
			word := p.currentQuery.String()[queryLen-patternLen:]
			patternFound = patternFound ||
				strings.ToUpper(word) == pattern
		}
	}

	return patternFound
}

func (p *Parser) Parse() []string {
	var commands []string

	for i := 0; i < len(p.Content); i++ {
		p.currentChar = p.Content[i]
		p.currentQuery.WriteRune(rune(p.currentChar))

		if p.isQueryComplete() {
			commands = append(commands, p.currentQuery.String())
			p.currentQuery.Reset()
			continue
		}

		p.handleLiterals()
		p.handleComments()
		p.handleTransactions()
	}

	return commands
}
