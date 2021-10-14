// Copyright (c) 2015 Taco de Wolff.
// Use of this source code is governed by a MIT style
// license that can be found in https://github.com/tdewolff/parse/blob/master/LICENSE.md
// Copyright (c) 2021 seasonjs Core Team
// Use of this source code is governed by a MIT style
// license that can be found in https://github.com/seasonjs/espack/blob/main/LICENSE

// Package js is an ECMAScript5.1 lexer following the specifications at http://www.ecma-international.org/ecma-262/5.1/.

package lexer

import (
	"github.com/seasonjs/espack/internal/builder/pkg/core/in"
	"unicode"
	"unicode/utf8"
)

var identifierStart = []*unicode.RangeTable{unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo, unicode.Nl, unicode.Other_ID_Start}
var identifierContinue = []*unicode.RangeTable{unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo, unicode.Nl, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc, unicode.Other_ID_Continue}

// IsIdentifierStart returns true if the byte-slice start is the start of an identifier
func IsIdentifierStart(b []byte) bool {
	r, _ := utf8.DecodeRune(b)
	return r == '$' || r == '\\' || r == '_' || unicode.IsOneOf(identifierStart, r)
}

// IsIdentifierContinue returns true if the byte-slice start is a continuation of an identifier
func IsIdentifierContinue(b []byte) bool {
	r, _ := utf8.DecodeRune(b)
	return r == '$' || r == '\\' || r == '\u200C' || r == '\u200D' || unicode.IsOneOf(identifierContinue, r)
}

// IsIdentifierEnd returns true if the byte-slice end is a start or continuation of an identifier
func IsIdentifierEnd(b []byte) bool {
	r, _ := utf8.DecodeLastRune(b)
	return r == '$' || r == '\\' || r == '\u200C' || r == '\u200D' || unicode.IsOneOf(identifierContinue, r)
}

////////////////////////////////////////////////////////////////

// Position 位置
type Position struct {
	line   int // >= 1
	column int // >= 0
}

type Location struct {
	startLoc Position
	endLoc   Position
}

// Lexer is the state for the lexer.
type Lexer struct {
	r                  *in.Input
	err                error
	prevLineTerminator bool
	prevNumericLiteral bool
	level              int
	templateLevels     []int
	loc                Location
	pos                int
	curLine            int
}

// NewLexer returns a new Lexer for a given io.Reader.
func NewLexer(r *in.Input) *Lexer {
	return &Lexer{
		r:                  r,
		prevLineTerminator: true,
		level:              0,
		templateLevels:     []int{},
		pos:                0,
		curLine:            1,
	}
}

// Err returns the error encountered during lexing, this is often io.EOF but also other errors can be returned.
func (l *Lexer) Err() error {
	if l.err != nil {
		return l.err
	}
	return l.r.Err()
}

// RegExp reparses the input stream for a regular expression. It is assumed that we just received DivToken or DivEqToken with Next(). This function will go back and read that as a regular expression.
func (l *Lexer) RegExp() (TokenType, []byte) {
	if 0 < l.r.Offset() && l.r.Peek(-1) == '/' {
		l.r.Move(-1)
	} else if 1 < l.r.Offset() && l.r.Peek(-1) == '=' && l.r.Peek(-2) == '/' {
		l.r.Move(-2)
	} else {
		l.err = in.NewErrorLexer(l.r, "expected / or /=")
		return ErrorToken, nil
	}
	l.r.Skip() // trick to set start = pos

	if l.consumeRegExpToken() {
		return RegExpToken, l.r.Shift()
	}
	l.err = in.NewErrorLexer(l.r, "unexpected EOF or newline")
	return ErrorToken, nil
}

// Next 返回一个Token状态信息
func (l *Lexer) Next() (TokenType, []byte, Location) {
	prevLineTerminator := l.prevLineTerminator
	l.prevLineTerminator = false

	prevNumericLiteral := l.prevNumericLiteral
	l.prevNumericLiteral = false

	// study on 50x jQuery shows:
	// spaces: 20k
	// alpha: 16k
	// newlines: 14.4k
	// operators: 4k
	// numbers and dot: 3.6k
	// (): 3.4k
	// {}: 1.8k
	// []: 0.9k
	// "': 1k
	// semicolon: 2.4k
	// colon: 0.8k
	// comma: 2.4k
	// slash: 1.4k
	// `~: almost 0

	c := l.r.Peek(0)
	switch c {
	case ' ', '\t', '\v', '\f':
		start := l.curPosition()
		l.r.Move(1)
		for l.consumeWhitespace() {
		}
		end := l.curPosition()
		l.loc = Location{
			start, end,
		}
		l.prevLineTerminator = prevLineTerminator
		return WhitespaceToken, l.r.Shift(), l.loc
	case '\n', '\r':
		l.curLine++
		start := l.curPosition()
		l.r.Move(1)
		for l.consumeLineTerminator() {
			l.curLine++
		}
		end := l.curPosition()
		l.pos = 0
		l.loc = Location{
			start, end,
		}
		l.prevLineTerminator = true
		return LineTerminatorToken, l.r.Shift(), l.loc
	case '>', '=', '!', '+', '*', '%', '&', '|', '^', '~', '?':
		start := l.curPosition()

		if tt := l.consumeOperatorToken(); tt != ErrorToken {
			end := l.curPosition()
			l.loc = Location{
				start, end,
			}
			return tt, l.r.Shift(), l.loc
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
		start := l.curPosition()

		if tt := l.consumeNumericToken(); tt != ErrorToken || l.r.Pos() != 0 {
			end := l.curPosition()
			l.loc = Location{
				start, end,
			}
			l.prevNumericLiteral = true
			return tt, l.r.Shift(), l.loc
		} else if c == '.' {
			l.r.Move(1)
			if l.r.Peek(0) == '.' && l.r.Peek(1) == '.' {
				l.r.Move(2)
				end := l.curPosition()
				l.loc = Location{
					start, end,
				}
				return EllipsisToken, l.r.Shift(), l.loc
			}
			end := l.curPosition()
			l.loc = Location{
				start, end,
			}
			return DotToken, l.r.Shift(), l.loc
		}
	case ',':
		start := l.curPosition()
		l.r.Move(1)
		end := l.curPosition()
		l.loc = Location{
			start, end,
		}
		return CommaToken, l.r.Shift(), l.loc
	case ';':
		start := l.curPosition()
		l.r.Move(1)
		end := l.curPosition()
		l.loc = Location{
			start, end,
		}
		return SemicolonToken, l.r.Shift(), l.loc
	case '(':
		start := l.curPosition()
		l.level++
		l.r.Move(1)
		end := l.curPosition()
		l.loc = Location{
			start, end,
		}
		return OpenParenToken, l.r.Shift(), l.loc
	case ')':
		l.level--
		l.r.Move(1)
		return CloseParenToken, l.r.Shift(), l.loc
	case '/':
		start := l.curPosition()
		if tt := l.consumeCommentToken(); tt != ErrorToken {
			end := l.curPosition()
			l.loc = Location{
				start, end,
			}
			return tt, l.r.Shift(), l.loc
		} else if tt := l.consumeOperatorToken(); tt != ErrorToken {
			end := l.curPosition()
			l.loc = Location{
				start, end,
			}
			return tt, l.r.Shift(), l.loc
		}
	case '{':
		start := l.curPosition()
		l.level++
		l.r.Move(1)
		end := l.curPosition()
		l.loc = Location{
			start, end,
		}
		return OpenBraceToken, l.r.Shift(), l.loc
	case '}':
		start := l.curPosition()
		l.level--
		if len(l.templateLevels) != 0 && l.level == l.templateLevels[len(l.templateLevels)-1] {
			tt, end := l.consumeTemplateToken()
			l.loc = Location{
				start, end,
			}
			return tt, l.r.Shift(), l.loc
		}
		l.r.Move(1)
		end := l.curPosition()
		l.loc = Location{
			start, end,
		}
		return CloseBraceToken, l.r.Shift(), l.loc
	case ':':
		start := l.curPosition()
		l.r.Move(1)
		l.r.Move(1)
		end := l.curPosition()
		l.loc = Location{
			start, end,
		}
		return ColonToken, l.r.Shift(), l.loc
	case '\'', '"':
		start := l.curPosition()
		if l.consumeStringToken() {
			end := l.curPosition()
			l.loc = Location{
				start, end,
			}
			return StringToken, l.r.Shift(), l.loc
		}
	case ']':
		start := l.curPosition()
		l.r.Move(1)
		end := l.curPosition()
		l.loc = Location{
			start, end,
		}
		return CloseBracketToken, l.r.Shift(), l.loc
	case '[':
		start := l.curPosition()
		l.r.Move(1)
		end := l.curPosition()
		l.loc = Location{
			start, end,
		}
		return OpenBracketToken, l.r.Shift(), l.loc
	case '<', '-':
		start := l.curPosition()
		if l.consumeHTMLLikeCommentToken(prevLineTerminator) {
			end := l.curPosition()
			l.loc = Location{
				start, end,
			}
			return CommentToken, l.r.Shift(), l.loc
		} else if tt := l.consumeOperatorToken(); tt != ErrorToken {
			end := l.curPosition()
			l.loc = Location{
				start, end,
			}
			return tt, l.r.Shift(), l.loc
		}
	case '`':
		start := l.curPosition()
		l.templateLevels = append(l.templateLevels, l.level)
		tt, end := l.consumeTemplateToken()
		l.loc = Location{
			start, end,
		}
		return tt, l.r.Shift(), l.loc
	case '#':
		start := l.curPosition()
		l.r.Move(1)
		if l.consumeIdentifierToken() {
			end := l.curPosition()
			l.loc = Location{
				start, end,
			}
			return PrivateIdentifierToken, l.r.Shift(), l.loc
		}
		end := l.curPosition()
		l.loc = Location{
			start, end,
		}
		return ErrorToken, nil, l.loc
	default:
		start := l.curPosition()
		if l.consumeIdentifierToken() {
			if prevNumericLiteral {
				end := l.curPosition()
				l.loc = Location{
					start, end,
				}
				l.err = in.NewErrorLexer(l.r, "unexpected identifier after number")
				return ErrorToken, nil, l.loc
			} else if keyword, ok := Keywords[string(l.r.Lexeme())]; ok {
				end := l.curPosition()
				l.loc = Location{
					start, end,
				}
				return keyword, l.r.Shift(), l.loc
			}
			end := l.curPosition()
			l.loc = Location{
				start, end,
			}
			return IdentifierToken, l.r.Shift(), l.loc
		}
		if 0xC0 <= c {
			if l.consumeWhitespace() {
				for l.consumeWhitespace() {
				}
				end := l.curPosition()
				l.loc = Location{
					start, end,
				}
				l.prevLineTerminator = prevLineTerminator
				return WhitespaceToken, l.r.Shift(), l.loc
			} else if l.consumeLineTerminator() {
				for l.consumeLineTerminator() {
				}
				end := l.curPosition()
				l.loc = Location{
					start, end,
				}
				l.prevLineTerminator = true
				return LineTerminatorToken, l.r.Shift(), l.loc
			}
		} else if c == 0 && l.r.Err() != nil {
			end := l.curPosition()
			l.loc = Location{
				start, end,
			}
			return ErrorToken, nil, l.loc
		}
	}

	r, _ := l.r.PeekRune(0)
	l.err = in.NewErrorLexer(l.r, "unexpected %s", in.Printable(r))
	return ErrorToken, l.r.Shift(), l.loc
}

////////////////////////////////////////////////////////////////

/*
The following functions follow the specifications at http://www.ecma-international.org/ecma-262/5.1/
*/

func (l *Lexer) curPosition() Position {
	//pos 存在问题 ：）它取的不是正确的line
	l.pos = l.pos + l.r.Pos()
	return Position{
		l.curLine, l.pos,
	}
}

func (l *Lexer) consumeWhitespace() bool {
	c := l.r.Peek(0)
	if c == ' ' || c == '\t' || c == '\v' || c == '\f' {
		l.r.Move(1)
		return true
	} else if 0xC0 <= c {
		if r, n := l.r.PeekRune(0); r == '\u00A0' || r == '\uFEFF' || unicode.Is(unicode.Zs, r) {
			l.r.Move(n)
			return true
		}
	}
	return false
}

func (l *Lexer) isLineTerminator() bool {
	c := l.r.Peek(0)
	if c == '\n' || c == '\r' {
		return true
	} else if c == 0xE2 && l.r.Peek(1) == 0x80 && (l.r.Peek(2) == 0xA8 || l.r.Peek(2) == 0xA9) {
		return true
	}
	return false
}

func (l *Lexer) consumeLineTerminator() bool {
	c := l.r.Peek(0)
	if c == '\n' {
		l.r.Move(1)
		return true
	} else if c == '\r' {
		if l.r.Peek(1) == '\n' {
			l.r.Move(2)
		} else {
			l.r.Move(1)
		}
		return true
	} else if c == 0xE2 && l.r.Peek(1) == 0x80 && (l.r.Peek(2) == 0xA8 || l.r.Peek(2) == 0xA9) {
		l.r.Move(3)
		return true
	}
	return false
}

func (l *Lexer) consumeDigit() bool {
	if c := l.r.Peek(0); c >= '0' && c <= '9' {
		l.r.Move(1)
		return true
	}
	return false
}

func (l *Lexer) consumeHexDigit() bool {
	if c := l.r.Peek(0); (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') {
		l.r.Move(1)
		return true
	}
	return false
}

func (l *Lexer) consumeBinaryDigit() bool {
	if c := l.r.Peek(0); c == '0' || c == '1' {
		l.r.Move(1)
		return true
	}
	return false
}

func (l *Lexer) consumeOctalDigit() bool {
	if c := l.r.Peek(0); c >= '0' && c <= '7' {
		l.r.Move(1)
		return true
	}
	return false
}

func (l *Lexer) consumeUnicodeEscape() bool {
	if l.r.Peek(0) != '\\' || l.r.Peek(1) != 'u' {
		return false
	}
	mark := l.r.Pos()
	l.r.Move(2)
	if c := l.r.Peek(0); c == '{' {
		l.r.Move(1)
		if l.consumeHexDigit() {
			for l.consumeHexDigit() {
			}
			if c := l.r.Peek(0); c == '}' {
				l.r.Move(1)
				return true
			}
		}
		l.r.Rewind(mark)
		return false
	} else if !l.consumeHexDigit() || !l.consumeHexDigit() || !l.consumeHexDigit() || !l.consumeHexDigit() {
		l.r.Rewind(mark)
		return false
	}
	return true
}

func (l *Lexer) consumeSingleLineComment() {
	for {
		c := l.r.Peek(0)
		if c == '\r' || c == '\n' || c == 0 && l.r.Err() != nil {
			break
		} else if 0xC0 <= c {
			if r, _ := l.r.PeekRune(0); r == '\u2028' || r == '\u2029' {
				break
			}
		}
		l.r.Move(1)
	}
}

////////////////////////////////////////////////////////////////

func (l *Lexer) consumeHTMLLikeCommentToken(prevLineTerminator bool) bool {
	c := l.r.Peek(0)
	if c == '<' && l.r.Peek(1) == '!' && l.r.Peek(2) == '-' && l.r.Peek(3) == '-' {
		// opening HTML-style single line comment
		l.r.Move(4)
		l.consumeSingleLineComment()
		return true
	} else if prevLineTerminator && c == '-' && l.r.Peek(1) == '-' && l.r.Peek(2) == '>' {
		// closing HTML-style single line comment
		// (only if current line didn't contain any meaningful tokens)
		l.r.Move(3)
		l.consumeSingleLineComment()
		return true
	}
	return false
}

func (l *Lexer) consumeCommentToken() TokenType {
	c := l.r.Peek(1)
	if c == '/' {
		// single line comment
		l.r.Move(2)
		l.consumeSingleLineComment()
		return CommentToken
	} else if c == '*' {
		l.r.Move(2)
		tt := CommentToken
		for {
			c := l.r.Peek(0)
			if c == '*' && l.r.Peek(1) == '/' {
				l.r.Move(2)
				break
			} else if c == 0 && l.r.Err() != nil {
				break
			} else if l.consumeLineTerminator() {
				l.prevLineTerminator = true
				tt = CommentLineTerminatorToken
			} else {
				l.r.Move(1)
			}
		}
		return tt
	}
	return ErrorToken
}

var opTokens = map[byte]TokenType{
	'=': EqToken,
	'!': NotToken,
	'<': LtToken,
	'>': GtToken,
	'+': AddToken,
	'-': SubToken,
	'*': MulToken,
	'/': DivToken,
	'%': ModToken,
	'&': BitAndToken,
	'|': BitOrToken,
	'^': BitXorToken,
	'~': BitNotToken,
	'?': QuestionToken,
}

var opEqTokens = map[byte]TokenType{
	'=': EqEqToken,
	'!': NotEqToken,
	'<': LtEqToken,
	'>': GtEqToken,
	'+': AddEqToken,
	'-': SubEqToken,
	'*': MulEqToken,
	'/': DivEqToken,
	'%': ModEqToken,
	'&': BitAndEqToken,
	'|': BitOrEqToken,
	'^': BitXorEqToken,
}

var opOpTokens = map[byte]TokenType{
	'<': LtLtToken,
	'+': IncrToken,
	'-': DecrToken,
	'*': ExpToken,
	'&': AndToken,
	'|': OrToken,
	'?': NullishToken,
}

var opOpEqTokens = map[byte]TokenType{
	'<': LtLtEqToken,
	'*': ExpEqToken,
	'&': AndEqToken,
	'|': OrEqToken,
	'?': NullishEqToken,
}

func (l *Lexer) consumeOperatorToken() TokenType {
	c := l.r.Peek(0)
	l.r.Move(1)
	if l.r.Peek(0) == '=' {
		l.r.Move(1)
		if l.r.Peek(0) == '=' && (c == '!' || c == '=') {
			l.r.Move(1)
			if c == '!' {
				return NotEqEqToken
			}
			return EqEqEqToken
		}
		return opEqTokens[c]
	} else if l.r.Peek(0) == c && (c == '+' || c == '-' || c == '*' || c == '&' || c == '|' || c == '?' || c == '<') {
		l.r.Move(1)
		if l.r.Peek(0) == '=' && c != '+' && c != '-' {
			l.r.Move(1)
			return opOpEqTokens[c]
		}
		return opOpTokens[c]
	} else if c == '?' && l.r.Peek(0) == '.' && (l.r.Peek(1) < '0' || l.r.Peek(1) > '9') {
		l.r.Move(1)
		return OptChainToken
	} else if c == '=' && l.r.Peek(0) == '>' {
		l.r.Move(1)
		return ArrowToken
	} else if c == '>' && l.r.Peek(0) == '>' {
		l.r.Move(1)
		if l.r.Peek(0) == '>' {
			l.r.Move(1)
			if l.r.Peek(0) == '=' {
				l.r.Move(1)
				return GtGtGtEqToken
			}
			return GtGtGtToken
		} else if l.r.Peek(0) == '=' {
			l.r.Move(1)
			return GtGtEqToken
		}
		return GtGtToken
	}
	return opTokens[c]
}

func (l *Lexer) consumeIdentifierToken() bool {
	c := l.r.Peek(0)
	if identifierStartTable[c] {
		l.r.Move(1)
	} else if 0xC0 <= c {
		if r, n := l.r.PeekRune(0); unicode.IsOneOf(identifierStart, r) {
			l.r.Move(n)
		} else {
			return false
		}
	} else if !l.consumeUnicodeEscape() {
		return false
	}
	for {
		c := l.r.Peek(0)
		if identifierTable[c] {
			l.r.Move(1)
		} else if 0xC0 <= c {
			if r, n := l.r.PeekRune(0); r == '\u200C' || r == '\u200D' || unicode.IsOneOf(identifierContinue, r) {
				l.r.Move(n)
			} else {
				break
			}
		} else {
			break
		}
	}
	return true
}

func (l *Lexer) consumeNumericToken() TokenType {
	// assume to be on 0 1 2 3 4 5 6 7 8 9 .
	first := l.r.Peek(0)
	if first == '0' {
		l.r.Move(1)
		if l.r.Peek(0) == 'x' || l.r.Peek(0) == 'X' {
			l.r.Move(1)
			if l.consumeHexDigit() {
				for l.consumeHexDigit() {
				}
				return HexadecimalToken
			}
			l.err = in.NewErrorLexer(l.r, "invalid hexadecimal number")
			return ErrorToken
		} else if l.r.Peek(0) == 'b' || l.r.Peek(0) == 'B' {
			l.r.Move(1)
			if l.consumeBinaryDigit() {
				for l.consumeBinaryDigit() {
				}
				return BinaryToken
			}
			l.err = in.NewErrorLexer(l.r, "invalid binary number")
			return ErrorToken
		} else if l.r.Peek(0) == 'o' || l.r.Peek(0) == 'O' {
			l.r.Move(1)
			if l.consumeOctalDigit() {
				for l.consumeOctalDigit() {
				}
				return OctalToken
			}
			l.err = in.NewErrorLexer(l.r, "invalid octal number")
			return ErrorToken
		} else if l.r.Peek(0) == 'n' {
			l.r.Move(1)
			return BigIntToken
		} else if '0' <= l.r.Peek(0) && l.r.Peek(0) <= '9' {
			l.err = in.NewErrorLexer(l.r, "legacy octal numbers are not supported")
			return ErrorToken
		}
	} else if first != '.' {
		for l.consumeDigit() {
		}
	}
	// we have parsed a 0 or an integer number
	c := l.r.Peek(0)
	if c == '.' {
		l.r.Move(1)
		if l.consumeDigit() {
			for l.consumeDigit() {
			}
			c = l.r.Peek(0)
		} else if first == '.' {
			// number starts with a dot and must be followed by digits
			l.r.Move(-1)
			return ErrorToken // may be dot or ellipsis
		} else {
			c = l.r.Peek(0)
		}
	} else if c == 'n' {
		l.r.Move(1)
		return BigIntToken
	}
	if c == 'e' || c == 'E' {
		l.r.Move(1)
		c = l.r.Peek(0)
		if c == '+' || c == '-' {
			l.r.Move(1)
		}
		if !l.consumeDigit() {
			l.err = in.NewErrorLexer(l.r, "invalid number")
			return ErrorToken
		}
		for l.consumeDigit() {
		}
	}
	return DecimalToken
}

func (l *Lexer) consumeStringToken() bool {
	// assume to be on ' or "
	mark := l.r.Pos()
	delim := l.r.Peek(0)
	l.r.Move(1)
	for {
		c := l.r.Peek(0)
		if c == delim {
			l.r.Move(1)
			break
		} else if c == '\\' {
			l.r.Move(1)
			if !l.consumeLineTerminator() {
				if c := l.r.Peek(0); c == delim || c == '\\' {
					l.r.Move(1)
				}
			}
			continue
		} else if c == '\n' || c == '\r' || c == 0 && l.r.Err() != nil {
			l.r.Rewind(mark)
			return false
		}
		l.r.Move(1)
	}
	return true
}

func (l *Lexer) consumeRegExpToken() bool {
	// assume to be on /
	l.r.Move(1)
	inClass := false
	for {
		c := l.r.Peek(0)
		if !inClass && c == '/' {
			l.r.Move(1)
			break
		} else if c == '[' {
			inClass = true
		} else if c == ']' {
			inClass = false
		} else if c == '\\' {
			l.r.Move(1)
			if l.isLineTerminator() || l.r.Peek(0) == 0 && l.r.Err() != nil {
				return false
			}
		} else if l.isLineTerminator() || c == 0 && l.r.Err() != nil {
			return false
		}
		l.r.Move(1)
	}
	// flags
	for {
		c := l.r.Peek(0)
		if identifierTable[c] {
			l.r.Move(1)
		} else if 0xC0 <= c {
			if r, n := l.r.PeekRune(0); r == '\u200C' || r == '\u200D' || unicode.IsOneOf(identifierContinue, r) {
				l.r.Move(n)
			} else {
				break
			}
		} else {
			break
		}
	}
	return true
}

func (l *Lexer) consumeTemplateToken() (TokenType, Position) {
	// assume to be on ` or } when already within template
	continuation := l.r.Peek(0) == '}'
	l.r.Move(1)
	for {
		c := l.r.Peek(0)
		if c == '`' {
			l.templateLevels = l.templateLevels[:len(l.templateLevels)-1]
			l.r.Move(1)
			if continuation {
				return TemplateEndToken, l.curPosition()
			}
			return TemplateToken, l.curPosition()
		} else if c == '$' && l.r.Peek(1) == '{' {
			l.level++
			l.r.Move(2)
			if continuation {
				return TemplateMiddleToken, l.curPosition()
			}
			return TemplateStartToken, l.curPosition()
		} else if c == '\\' {
			l.r.Move(1)
			if c := l.r.Peek(0); c != 0 {
				l.r.Move(1)
			}
			continue
		} else if c == 0 && l.r.Err() != nil {
			if continuation {
				return TemplateEndToken, l.curPosition()
			}
			return TemplateToken, l.curPosition()
		}
		l.r.Move(1)
	}
}

var identifierStartTable = [256]bool{
	// ASCII
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, true, false, false, false, // $
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, true, true, true, true, true, true, true, // A, B, C, D, E, F, G
	true, true, true, true, true, true, true, true, // H, I, J, K, L, M, N, O
	true, true, true, true, true, true, true, true, // P, Q, R, S, T, U, V, W
	true, true, true, false, false, false, false, true, // X, Y, Z, _

	false, true, true, true, true, true, true, true, // a, b, c, d, e, f, g
	true, true, true, true, true, true, true, true, // h, i, j, k, l, m, n, o
	true, true, true, true, true, true, true, true, // p, q, r, s, t, u, v, w
	true, true, true, false, false, false, false, false, // x, y, z

	// non-ASCII
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
}

var identifierTable = [256]bool{
	// ASCII
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, true, false, false, false, // $
	false, false, false, false, false, false, false, false,
	true, true, true, true, true, true, true, true, // 0, 1, 2, 3, 4, 5, 6, 7
	true, true, false, false, false, false, false, false, // 8, 9

	false, true, true, true, true, true, true, true, // A, B, C, D, E, F, G
	true, true, true, true, true, true, true, true, // H, I, J, K, L, M, N, O
	true, true, true, true, true, true, true, true, // P, Q, R, S, T, U, V, W
	true, true, true, false, false, false, false, true, // X, Y, Z, _

	false, true, true, true, true, true, true, true, // a, b, c, d, e, f, g
	true, true, true, true, true, true, true, true, // h, i, j, k, l, m, n, o
	true, true, true, true, true, true, true, true, // p, q, r, s, t, u, v, w
	true, true, true, false, false, false, false, false, // x, y, z

	// non-ASCII
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,

	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
}
