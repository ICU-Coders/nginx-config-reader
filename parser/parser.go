package parser

import (
	"github.com/fatih/color"
	"io/ioutil"
	"path/filepath"
	"strings"
	"unicode"
)

var ObjTokens = [...]string{"server", "http", "location", "events", "types"}

var PropertyTokens = [...]string{"root", "index", "alias", "proxy_pass", "listen", "error_page", "log_format",
	"access_log", "include", "user", "worker_processes", "error_log", "pid", "sendfile",
	"tcp_nopush", "tcp_nodelay", "keepalive_timeout", "types_hash_max_size", "worker_connections", "server_name"}

func tokenizer(input []byte) []Token {
	current := 0
	tokens := make([]Token, 0, 500)
	for current < len(input) {
		char := input[current]
		if char == '{' {
			tokens = append(tokens, Token{
				name:  TokenNameBrace,
				value: "{",
			})
			current++
			continue
		}

		if char == '}' {
			tokens = append(tokens, Token{
				name:  TokenNameBrace,
				value: "}",
			})
			current++
			continue
		}

		if char == ';' {
			tokens = append(tokens, Token{
				name:  TokenNameSemicolon,
				value: ";",
			})
			current++
			continue
		}

		if char == '#' {
			tokens = append(tokens, Token{
				name:  TokenNameComment,
				value: "#",
			})
			current++
			continue
		}

		if char == '\n' {
			tokens = append(tokens, Token{
				name:  TokenNameNewLine,
				value: "\n",
			})
			current++
			continue
		}

		if unicode.IsSpace(rune(char)) {
			current++
			continue
		}

		{
			var value []byte
			for !unicode.IsSpace(rune(char)) && char != '#' && char != ';' {
				value = append(value, char)
				current++
				char = input[current]
			}
			str := string(value)
			find := false
			for _, token := range ObjTokens {
				if token == str {
					find = true
					break
				}
			}
			if !find {
				for _, token := range PropertyTokens {
					if token == str {
						find = true
						break
					}
				}
			}
			if !find {
				tokens = append(tokens, Token{
					name:  TokenNameString,
					value: str,
				})
			} else {
				tokens = append(tokens, Token{
					name:  TokenNameToken,
					value: str,
				})
			}
		}
	}
	return tokens
}

func parser(tokens []Token) *ObjectExpression {
	current := 0
	var walk func(end TokenName) Node
	walk = func(end TokenName) Node {
		token := tokens[current]
		switch end {
		case TokenNameNewLine, TokenNameSemicolon, TokenNameBrace:
			m := new(MultiValue)
			for token.name != end {
				m.Values = append(m.Values, token.value)
				current++
				token = tokens[current]
			}
			return m
		}

		if token.name == TokenNameComment {
			current++
			return &PropertyExpression{Key: "common", Value: walk(TokenNameNewLine)}
		}
		if token.name == TokenNameToken {
			p := new(PropertyExpression)
			p.Key = token.value

			find := false
			for _, objToken := range ObjTokens {
				if objToken == token.value {
					find = true
					break
				}
			}
			current++
			if !find {
				p.Value = walk(TokenNameSemicolon).(*MultiValue)
			} else {
				if token.value == "location" {
					path := walk(TokenNameBrace)
					body := walk(TokenNameAny)
					p.Value = &ObjectExpression{
						Properties: []*PropertyExpression{
							{
								Key:   "path",
								Value: path,
							},
							{
								Key:   "body",
								Value: body,
							},
						},
					}
				} else {
					p.Value = walk(TokenNameAny)
				}
			}
			return p
		}

		if token.name == TokenNameBrace && token.value == "{" {
			obj := new(ObjectExpression)
			current++
			token = tokens[current]

			for token.name != TokenNameBrace || (token.name == TokenNameBrace && token.value != "}") {
				obj.Append(walk(TokenNameAny))
				token = tokens[current]
			}
			current++
			return obj
		}

		// not match
		current++
		return &PropertyExpression{Key: "", Value: nil} // empty
	}

	body := new(ObjectExpression)
	for current < len(tokens) {
		body.Append(walk(TokenNameAny))
	}

	return body
}

func ConfigParser(path string) (*ObjectExpression, error) {
	path = strings.TrimSpace(path)
	color.New(color.FgCyan).Println("> Parsing", path)
	if isVisited(path) {
		return nil, nil
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tokens := tokenizer(file)
	nginx := parser(tokens)
	// 遍历展开include
	err = extraNginx(nginx)
	if err != nil {
		return nil, err
	}
	return nginx, nil
}

func extraNginx(nginx *ObjectExpression) error {
	includes := nginx.GetMustString("include")
	for _, include := range includes {
		glob, err := filepath.Glob(include)
		if err != nil {
			return err
		}
		for _, s := range glob {
			obj, err := ConfigParser(s)
			if err != nil {
				return err
			}
			if obj == nil {
				continue
			}
			nginx.Properties = append(nginx.Properties, obj.Properties...)
		}
	}

	https := nginx.GetMustObject("http")
	for _, http := range https {
		err := extraNginx(http)
		if err != nil {
			return err
		}
	}

	servers := nginx.GetMustObject("server")
	for _, server := range servers {
		err := extraNginx(server)
		if err != nil {
			return err
		}
	}

	return nil
}
