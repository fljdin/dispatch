package parser

import "strings"

type ShParser struct {
	content string
}

func (p *ShParser) Parse() []string {
	lines := strings.Split(p.content, `\n`)
	commands := make([]string, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			commands = append(commands, line)
		}
	}

	return commands
}

func (p *ShParser) SetContent(content string) {
	p.content = content
}
