package pacgo

import (
	"errors"
	"strings"
)

// Parser is an empty struct user to parse PKGBUILD content
type Parser struct{}

// NewParser returns a new Parser object
func NewParser() *Parser {
	return &Parser{}
}

// Parse method parses the content passed (as a PKGBUILD) and returns a pointer to
// a PkgBuild object and an error if any.
func (p *Parser) Parse(content string) (*PkgBuild, error) {
	pinfo := NewPackageInfo()
	metadata, funcs := separateFuncs(abridge(tokenize(rmCommentAndEmptyLine([]byte(content)))))
	if len(metadata)%3 != 0 {
		return nil, errors.New("Syntax Error")
	}
	for i := 0; i < len(metadata)/3; i++ {
		key := metadata[i*3]
		val := metadata[(i*3)+2]
		err := pinfo.Set(key, val)
		if err != nil {
			return nil, err
		}
	}
	return &PkgBuild{PackageInfo: pinfo, funcs: funcs}, nil
}

// separateFuncs is used to separate function token strings from the metadata ones
func separateFuncs(tokens []string) ([]string, []string) {
	metadata := make([]string, 0)
	funcs := make([]string, 0)

	for _, tk := range tokens {
		if strings.Contains(tk, "(){") {
			funcs = append(funcs, tk)
			continue
		}

		metadata = append(metadata, tk)
	}
	return metadata, funcs
}

// abridge is used to join all the tokenized function parts
func abridge(tokens []string) []string {
	ntokens := make([]string, 0)
	for _, token := range tokens {
		if strings.HasPrefix(token, "(){") {
			if len(ntokens) > 0 {
				ntokens[len(ntokens)-1] += token
				continue
			}
		}

		ntokens = append(ntokens, token)
	}
	return ntokens
}

// tokenize tokenizes the given byte array and returns the tokens as a strnig array
func tokenize(bs []byte) []string {
	currentToken := ""
	inBrackets := false
	inSingleQuotes := false
	inDoubleQuotes := false
	inFunc := 0
	tokens := make([]string, 0)

	for _, ch := range bs {
		if !inBrackets && !inDoubleQuotes && !inSingleQuotes && inFunc == 0 {
			if strings.Contains("=: \n\t", string(ch)) {
				if currentToken != "" {
					if strings.HasPrefix(currentToken, "{") {
						if len(tokens) != 0 {
							tokens[len(tokens)-1] += currentToken
							currentToken = ""
							continue
						}
					}
					tokens = append(tokens, currentToken)
					if ch == '=' {
						tokens = append(tokens, "=")
					}
					currentToken = ""
					continue
				} else {
					if ch == '\n' {
						tokens = append(tokens, "")
					}
					continue
				}
			}
		}

		if !inBrackets && !inDoubleQuotes && !inSingleQuotes && ch == '{' {
			inFunc += 1
		}

		if !inBrackets && !inDoubleQuotes && !inSingleQuotes && ch == '}' {
			inFunc -= 1
		}

		if ch == '(' {
			if inFunc == 0 && !inDoubleQuotes && !inSingleQuotes {
				inBrackets = true
			}
		}

		if ch == ')' {
			if inFunc == 0 && !inDoubleQuotes && !inSingleQuotes {

				inBrackets = false
			}
		}

		if ch == '\'' {
			inSingleQuotes = !inSingleQuotes
		}

		if ch == '"' {
			inSingleQuotes = !inSingleQuotes
		}

		currentToken += string(ch)
	}
	return tokens
}

// rmCommentAndEmptyLine is used to remove any empty lines and commented out regions from
// the passed content.
func rmCommentAndEmptyLine(bs []byte) []byte {
	str := string(bs)
	nstr := ""
	for _, line := range strings.Split(str, "\n") {
		lx := strings.TrimSpace(line)
		if !strings.HasPrefix(lx, "#") && lx != "" {
			nstr += line + "\n"
		}
	}
	return []byte(nstr)
}
