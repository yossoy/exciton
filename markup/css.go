package markup

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/gorilla/css/scanner"
)

type cssParser struct {
	s                *scanner.Scanner
	t                *scanner.Token
	classPrefix      string
	depth            int
	buf              bytes.Buffer
	spaceDropped     bool
	globalscopeBlock []bool
	localNameMap     map[string]string
}

var (
	nestedAtRules = []string{
		"@media", "@supports", "@document", "@page", "@font-face",
		"@keyframes", "@viewport", "@counter-style", "@font-feature-values",
	}
)

func isNestedAtRule(atRule string) bool {
	for _, s := range nestedAtRules {
		if atRule == s {
			return true
		}
	}
	return false
}

func newCSSParser(css string, prefix string) *cssParser {
	return &cssParser{
		s:            scanner.New(css),
		classPrefix:  prefix,
		t:            nil,
		localNameMap: make(map[string]string),
	}
}

func (cp *cssParser) fetch() *scanner.Token {
	if cp.t == nil {
		cp.t = cp.s.Next()
	}
	return cp.t
}

func (cp *cssParser) drop() *scanner.Token {
	//	cp.spaceDropped = false
	cp.t = nil
	return cp.fetch()
}

func (cp *cssParser) dropSpace() *scanner.Token {
	if cp.spaceDropped {
		cp.drop()
	} else {
		cp.replaceString(" ")
	}
	cp.spaceDropped = true
	return cp.fetch()
}

func (cp *cssParser) pass() *scanner.Token {
	cp.spaceDropped = false
	t := cp.fetch()
	cp.buf.WriteString(t.Value)
	return cp.drop()
}

func (cp *cssParser) writeString(s string) {
	cp.buf.WriteString(s)
}

func (cp *cssParser) replaceString(s string) *scanner.Token {
	cp.spaceDropped = false
	cp.buf.WriteString(s)
	return cp.drop()
}

func (cp *cssParser) isError() bool {
	return cp.fetch().Type == scanner.TokenError
}

func (cp *cssParser) err() error {
	if cp.isError() {
		return fmt.Errorf("Invalid token: %s", cp.fetch().String())
	}
	return nil
}

func (cp *cssParser) isChar(c string) bool {
	t := cp.fetch()
	return t.Type == scanner.TokenChar && t.Value == c
}

func (cp *cssParser) isEOF() bool {
	return cp.fetch().Type == scanner.TokenEOF
}

func (cp *cssParser) isValid() bool {
	return !cp.isEOF() && !cp.isError()
}

func (cp *cssParser) isWS() bool {
	return cp.fetch().Type == scanner.TokenS
}

func (cp *cssParser) isCommnet() bool {
	return cp.fetch().Type == scanner.TokenComment
}

func (cp *cssParser) isCDOorCDC() bool {
	t := cp.fetch()
	return t.Type == scanner.TokenCDO || t.Type == scanner.TokenCDC
}

func (cp *cssParser) isIgnore() bool {
	return cp.isWS() || cp.isCommnet() || cp.isCDOorCDC()
}

func (cp *cssParser) parseBOM() (bool, error) {
	if cp.fetch().Type == scanner.TokenBOM {
		cp.drop()
		return true, nil
	}
	return false, cp.err()
}

func (cp *cssParser) Parse() (string, error) {
	if _, err := cp.parseBOM(); err != nil {
		return "", err
	}
	if err := cp.parseCSS(); err != nil {
		return "", err
	}
	return cp.buf.String(), nil
}

func (cp *cssParser) parseCSS() error {
	for cp.isValid() {
		switch {
		case cp.isIgnore():
			cp.drop()
		case cp.isChar("{") || cp.isChar("}"):
			return fmt.Errorf("Invalid %s", cp.fetch().String())
		default:
			if err := cp.parseRules(false); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cp *cssParser) parseRules(inNestBlock bool) error {
	inBlock := false
	if cp.isChar("{") {
		inBlock = true
		//cp.buf.WriteString("{")
		cp.depth++
		if !inNestBlock {
			cp.pass()
		} else {
			cp.drop()
		}
	}

loop:
	for cp.isValid() {
		switch {
		case cp.isIgnore():
			//TODO: drop all ws?
			cp.dropSpace()
		case cp.isChar("}"):
			if !inBlock {
				return fmt.Errorf("Unexpected }: %s", cp.fetch().String())
			}
			cp.depth--
			if !inNestBlock {
				cp.pass()
			}
			break loop
		default:
			err := cp.parseRule()
			if err != nil {
				return err
			}
		}
	}
	return cp.err()
}

func (cp *cssParser) parseRule() error {
	if cp.fetch().Type == scanner.TokenAtKeyword {
		return cp.parseAtRule()
	}
	return cp.parseQualifiedRule()
}

func (cp *cssParser) parseAtRule() error {
	rule := cp.fetch().Value
	cp.pass()

loop:
	for cp.isValid() {
		switch {
		case cp.isChar(";"):
			cp.pass()
			break loop
		case cp.isChar("{"):
			if isNestedAtRule(rule) {
				if err := cp.parseRules(false); err != nil {
					return err
				}
			} else {
				if err := cp.parseDeclarations(); err != nil {
					return err
				}
			}
			break loop
		default:
			if _, err := cp.parsePrelude(true); err != nil {
				return err
			}
		}
	}
	return cp.err()
}

func (cp *cssParser) parseQualifiedRule() error {
	var prelude []string
loop:
	for cp.isValid() {
		switch {
		case cp.isChar("{"):
			if prelude == nil {
				return fmt.Errorf("[%d:%d] Unexpected { character: %s", cp.fetch().Line, cp.fetch().Column, cp.fetch().Value)
			}
			if err := cp.parseDeclarations(); err != nil {
				return err
			}
			break loop
		default:
			p, err := cp.parsePrelude(false)
			if err != nil {
				return err
			}
			prelude = p
		}
	}
	return cp.err()
}

func (cp *cssParser) isInGlobalScopeBlock() bool {
	if len(cp.globalscopeBlock) > 0 {
		return cp.globalscopeBlock[len(cp.globalscopeBlock)-1]
	}
	return false
}

func (cp *cssParser) pushGlobalScopeBlock(isGlobal bool) {
	cp.globalscopeBlock = append(cp.globalscopeBlock, isGlobal)
}

func (cp *cssParser) popGlobalScopeBlock() bool {
	ret := cp.globalscopeBlock[len(cp.globalscopeBlock)-1]
	cp.globalscopeBlock = cp.globalscopeBlock[:len(cp.globalscopeBlock)-1]
	return ret
}

func (cp *cssParser) makeLocalName(v string) string {
	//t := cp.fetch()
	return fmt.Sprintf(".%s-%s", cp.classPrefix, v[1:])
}

func (cp *cssParser) parsePrelude(inAtRule bool) ([]string, error) {
	var selectors []string
	var sb strings.Builder
	nextIsPsuedoClass := false
	nextCanHasBlock := false
	currentBlockIsGlobal := cp.isInGlobalScopeBlock()
	isInGlobalLocalFunc := false

	flush := func() {
		if sb.Len() > 0 {
			s := sb.String()
			s2 := s
			if !currentBlockIsGlobal && strings.HasPrefix(s, ".") {
				if s3, ok := cp.localNameMap[s]; ok {
					s2 = s3
				} else {
					s2 = cp.makeLocalName(s)
					for _, v := range cp.localNameMap {
						if v == s2 {
							panic("already exist key: " + s)
						}
					}
					cp.localNameMap[s] = s2
				}
			}
			selectors = append(selectors, s2)
			cp.writeString(s2)
			sb.Reset()
		}
	}

loop:
	for cp.isValid() {
		switch {
		case cp.isIgnore():
			cp.dropSpace()
			continue
		case cp.isChar(","):
			flush()
			cp.pass()
		case cp.isChar(":"):
			flush()
			nextIsPsuedoClass = true
			cp.drop()
		case cp.isChar(";"):
			break loop
		case cp.isChar("{"):
			flush()
			if nextCanHasBlock {
				cp.pushGlobalScopeBlock(currentBlockIsGlobal)
				err := cp.parseRules(true)
				if err != nil {
					return nil, err
				}
			} else {
				break loop
			}
		case cp.isChar("}"):
			if nextCanHasBlock {
				cp.drop()
				currentBlockIsGlobal = cp.popGlobalScopeBlock()
			} else {
				return nil, fmt.Errorf("[%d:%d] Unexpected } character: %s", cp.fetch().Line, cp.fetch().Column, cp.fetch().Value)
			}
		case cp.fetch().Type == scanner.TokenFunction:
			v := cp.fetch().Value
			if nextIsPsuedoClass && v == "global(" || v == "local(" {
				cp.pushGlobalScopeBlock(currentBlockIsGlobal)
				currentBlockIsGlobal = v == "global("
				isInGlobalLocalFunc = true
				flush()
				cp.drop()
			} else {
				flush()
				cp.replaceString(":" + v)
			}
			nextIsPsuedoClass = false
		case cp.isChar(")"):
			if isInGlobalLocalFunc {
				currentBlockIsGlobal = cp.popGlobalScopeBlock()
				flush()
				cp.drop()
				isInGlobalLocalFunc = false
			} else {
				flush()
				cp.pass()
			}
		default:
			if nextIsPsuedoClass {
				flush()
				cn := cp.fetch().Value
				if cn == "global" || cn == "local" {
					nextCanHasBlock = true
					cp.drop()
					currentBlockIsGlobal = (cn == "global")
				} else {
					cp.replaceString(":" + cn)
				}
				nextIsPsuedoClass = false
			} else {
				sb.WriteString(cp.fetch().Value)
				cp.drop()
				nextCanHasBlock = false
			}
		}
	}
	flush()
	for i, v := range selectors {
		fmt.Printf("[%d] %q\n", i, v)
	}
	return selectors, cp.err()
}

func (cp *cssParser) parseDeclarations() error {
	if cp.isChar("{") {
		cp.pass()
	}

loop:
	for cp.isValid() {
		switch {
		case cp.isIgnore():
			cp.dropSpace()
		case cp.isChar("}"):
			cp.pass()
			break loop
		default:
			if err := cp.parseDeclaration(); err != nil {
				return err
			}
		}
	}
	return cp.err()
}

func (cp *cssParser) parseDeclaration() error {
loop:
	for cp.isValid() {
		switch {
		case cp.isIgnore():
			cp.dropSpace()
		case cp.isChar(":"):
			cp.pass()
		case cp.isChar(";"):
			cp.pass()
			break loop
		case cp.isChar("}"):
			break loop
		default:
			cp.pass()
		}
	}
	return cp.err()
}

func parseBlock(s *scanner.Scanner, term string) error {
	for {
		t := s.Next()
		switch t.Type {
		case scanner.TokenEOF:
			if term != "" {
				return errors.New("invalid terminate")
			}
			return nil
		case scanner.TokenChar:
			if t.Value == term {
				return nil
			}

		}
	}
}

func convertKlassCSS(css string, classPrefix string) (string, map[string]string, error) {
	p := newCSSParser(css, classPrefix)
	parsed, err := p.Parse()
	if err != nil {
		return "", nil, err
	}
	return parsed, p.localNameMap, nil
}
