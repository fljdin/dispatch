package parser

import (
	"fmt"
	"strings"
)

var ParserTypes = []string{"psql"}

type Parser struct {
	Type    string
	Content string

	currentChar    byte
	currentQuote   byte // either ', " or $
	currentComment byte // either - or *
	currentTag     string
	activeTag      strings.Builder
	currentQuery   strings.Builder
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

func (p *Parser) handleQuotedStrings() {
	if p.match("'", "\"") {
		if !p.inQuotedLiteral() {
			// enter into string
			p.currentQuote = p.currentChar
		} else if p.currentQuote == p.currentChar {
			// exit a string
			p.currentQuote = 0x0
		}
	}
}

func (p *Parser) handleDollarQuotedStrings() {
	if p.match("$") {
		if !p.inQuotedLiteral() {
			if p.activeTag.Len() > 0 {
				// first tag occurrence has been found
				// enter into string
				p.updateActiveTag()
				p.currentTag = p.activeTag.String()
				p.currentQuote = p.currentChar
				p.activeTag.Reset()
			} else {
				// initialize first tag with a dollar sign
				p.updateActiveTag()
			}

		} else if p.currentQuote == p.currentChar {
			if p.activeTag.Len() > 0 {
				p.updateActiveTag()
				if p.activeTag.String() == p.currentTag {
					// second tag occurrence has been found
					// exit a string
					p.currentQuote = 0x0
					p.activeTag.Reset()
				}
			} else {
				// initialize second tag with a dollar sign
				p.updateActiveTag()
			}
		}
	}

	if p.activeTag.Len() > 0 && p.isValidIdentifier(p.currentChar) {
		// construct active tag with any valid identifier
		p.updateActiveTag()
	}
}

func (p *Parser) updateActiveTag() {
	p.activeTag.WriteRune(rune(p.currentChar))
}

func (p *Parser) handleStrings() {
	p.handleQuotedStrings()
	p.handleDollarQuotedStrings()
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

func (p *Parser) isValidIdentifier(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
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

		p.handleStrings()
		p.handleComments()
		p.handleTransactions()
	}

	return commands
}
