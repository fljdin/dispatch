package parser

type Parser interface {
	Parse() []string
	SetContent(content string)
}
