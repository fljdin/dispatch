package parser

import (
	"bufio"
	"bytes"
	"strings"
)

type ShParser struct {
	content string
}

func (p *ShParser) SetContent(content string) {
	p.content = content
}

func (p *ShParser) Parse() []string {
	var commands []string

	scanner := bufio.NewScanner(strings.NewReader(p.content))
	scanner.Split(scanShellCommands)

	var command strings.Builder
	inComment := false

	for scanner.Scan() {
		line := scanner.Text()

		if len(strings.TrimSpace(line)) == 0 {
			// skip empty lines
			continue
		}

		if idx := strings.Index(line, "#"); idx >= 0 {
			// ignore comments
			line = line[:idx]
		}

		if strings.HasSuffix(line, "\\") && !inComment {
			// remove trailing backslash
			line = line[:len(line)-1]
			command.WriteString(line)
		} else {
			command.WriteString(line)
			commands = append(commands, strings.TrimSpace(command.String()))
			command.Reset()
		}
	}

	if command.Len() > 0 {
		// add the last command if it is not empty
		commands = append(commands, strings.TrimSpace(command.String()))
	}

	return commands
}

func scanShellCommands(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// return the command up to the newline character
		return i + 1, data[0:i], nil
	}

	if atEOF {
		// if we're at EOF, return the remaining data
		return len(data), data, nil
	}

	// request more data
	return 0, nil, nil
}
