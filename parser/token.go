package parser

import "fmt"

type TokenName int

const (
	TokenNameString TokenName = iota
	TokenNameToken
	TokenNameBrace
	TokenNameSemicolon
	TokenNameComment
	TokenNameNewLine
	TokenNameAny
)

type Token struct {
	name  TokenName
	value string
}

func (receiver Token) String() string {
	name := ""
	value := receiver.value
	switch receiver.name {
	case TokenNameNewLine:
		name = "new-line"
		value = "\\n"
	case TokenNameComment:
		name = "comment"
	case TokenNameSemicolon:
		name = "semicolon"
	case TokenNameToken:
		name = "token"
	case TokenNameBrace:
		name = "brace"
	case TokenNameString:
		name = "string"
	}
	return fmt.Sprintf("%s: %s", name, value)
}
