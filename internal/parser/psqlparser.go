package parser

import (
	"strings"

	"golang.org/x/exp/slices"
)

var SupportedCommands = []string{"g", "gdesc", "gexec", "gx", "crosstabview"}

type PsqlParser struct {
	content string

	currentChar    byte
	currentQuote   byte // either ', " or $
	currentComment byte // either - or *
	currentTag     string
	activeTag      strings.Builder
	activeCommand  strings.Builder
	currentQuery   strings.Builder
	inTransaction  bool
	inCommand      bool
}

func (p *PsqlParser) inComment() bool {
	return p.currentComment != 0x0
}

func (p *PsqlParser) inQuotedString() bool {
	return p.currentQuote != 0x0
}

func (p *PsqlParser) inCommentOrString() bool {
	return p.inComment() || p.inQuotedString()
}

func (p *PsqlParser) handleTransactions() {
	if p.match("BEGIN") && !p.inCommentOrString() {
		p.inTransaction = true
	}

	if p.match("END", "COMMIT", "ROLLBACK") && !p.inCommentOrString() {
		p.inTransaction = false
	}
}

func (p *PsqlParser) handleComments() {
	// opening a comment
	if p.match("--", "/*") && !p.inCommentOrString() {
		p.currentComment = p.currentChar
	}

	// ending a comment
	if (p.currentComment == '-' && p.currentChar == '\n') ||
		(p.currentComment == '*' && p.match("*/")) {
		p.currentComment = 0x0
	}
}

func (p *PsqlParser) handleStrings() {
	p.handleQuotedStrings()
	p.handleDollarQuotedStrings()
}

func (p *PsqlParser) handleQuotedStrings() {
	if p.match("'", "\"") {
		if !p.inCommentOrString() {
			// enter into string
			p.currentQuote = p.currentChar
		} else if p.currentQuote == p.currentChar {
			// exit a string
			p.currentQuote = 0x0
		}
	}
}

func (p *PsqlParser) handleDollarQuotedStrings() {
	if p.match("$") {
		if !p.inCommentOrString() {
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

func (p *PsqlParser) updateActiveTag() {
	p.activeTag.WriteRune(rune(p.currentChar))
}

func (p *PsqlParser) handleCommand() {
	if p.match("\\") && !p.inCommentOrString() {
		p.inCommand = true
		return
	}

	if p.inCommand {
		p.activeCommand.WriteRune(rune(p.currentChar))
	}
}

func (p *PsqlParser) isValidIdentifier(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
}

func (p *PsqlParser) retrieveCommand() string {
	c := p.activeCommand.String()
	r := strings.Fields(c)
	if len(r) > 0 {
		return r[0]
	}
	return ""
}

func (p *PsqlParser) isCommandComplete() bool {
	if p.currentChar == '\n' && p.inCommand {
		command := p.retrieveCommand()
		p.inCommand = false
		p.activeCommand.Reset()
		return slices.Contains(SupportedCommands, command)
	}
	return false
}

func (p *PsqlParser) isQueryComplete() bool {
	if p.inTransaction || p.inCommentOrString() {
		return false
	}

	return p.currentChar == ';' || p.isCommandComplete()
}

func (p *PsqlParser) match(pattern ...string) bool {
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

func (p *PsqlParser) SetContent(content string) {
	p.content = content
}

func (p *PsqlParser) Parse() []string {
	var commands []string

	for i := 0; i < len(p.content); i++ {
		p.currentChar = p.content[i]
		p.currentQuery.WriteRune(rune(p.currentChar))

		if p.isQueryComplete() {
			commands = append(commands, p.currentQuery.String())
			p.currentQuery.Reset()
			continue
		}

		p.handleStrings()
		p.handleCommand()
		p.handleComments()
		p.handleTransactions()
	}

	return commands
}
