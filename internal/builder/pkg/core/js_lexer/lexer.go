package js_lexer

import (
	"github.com/seasonjs/espack/internal/builder/pkg/core/types"
	"github.com/seasonjs/espack/internal/builder/pkg/core/util"
	"github.com/seasonjs/espack/internal/logger"
	"io"
	"io/ioutil"
)

// 前置知识了解：https://go.dev/blog/strings
// TODO 需要保证payload不会溢出

type Lexer struct {
	payload     []byte          //要扫描的代码
	index       int             //当前的位置
	len         int             // payload的总长度
	startLine   int             //当前token开始行数
	startColumn int             //当前token开始的列数
	endLine     int             //当前token结束的行数
	endColumn   int             //当前token结束的列数
	tokenType   types.TokenType //当前Token类型
	value       []byte          //当前Token的值
	err         error           //当前遇到的错误
}

func NewLexer(reader io.Reader) Lexer {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		logger.Error("scanner init error %v", err)
	}
	return Lexer{
		payload: buf,
		len:     len(buf),
	}
}

// isSpace 如果是空格
func isSpace(slice rune) bool {
	if slice <= '\u00FF' {
		// ASCII 码
		switch slice {
		case ' ', '\t', '\v', '\f':
			return true
		case '\u00A0': // 转译后的空格\u00A0是连续不断的空格
			return true
		}
		return false
	}
	// 	高位转译的空格
	if '\u2000' <= slice && slice <= '\u200a' {
		return true
	}
	switch slice {
	case '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}

// isLineTerminator 如果是换行符
func isLineTerminator(slice rune) bool {
	if slice <= '\u00FF' {
		// ASCII 码
		switch byte(slice) {
		case util.CarriageReturn, util.LineFeed:
			return true
		case '\u0085': // 转译后的'\u0085' 是下一行
			return true
		}
		return false
	}
	switch slice {
	case '\u2028', '\u2029':
		return true
	}
	return false
}

// Next 解析下一个Token，跳过空格和换行，然后自动增加index 和 line column
func (s Lexer) Next() {
	s.skipSpaceAndLineTerminator()
	s.startColumn = s.endColumn
	s.startLine = s.endLine
	ch := s.Peek(0)
	//判断是不是js的token
	switch ch {
	case util.Dot:
		s.consumeDot()
		return
	case util.LeftParenthesis:
		s.finishToken(types.OpenParenToken, 1)
		return
	case util.RightParenthesis:
		s.finishToken(types.CloseParenToken, 1)
		return
	case util.Semicolon:
		s.finishToken(types.SemicolonToken, 1)
		return
	case util.Comma:
		s.finishToken(types.CommaToken, 1)
		return
	case util.LeftSquareBracket:
		s.finishToken(types.OpenBracketToken, 1)
		return
	case util.RightSquareBracket:
		s.finishToken(types.CloseBracketToken, 1)
		return
	case util.LeftCurlyBrace:
		s.finishToken(types.OpenBraceToken, 1)
		return
	case util.RightCurlyBrace:
		s.finishToken(types.CloseBraceToken, 1)
		return
	case util.Colon:
		//TODO 考虑function bind处理
		s.finishToken(types.ColonToken, 1)
		return
	case util.QuotationMark:
		s.consumeQuotationMark()
		return
	case util.GraveAccent:
		s.consumeTemplate()
		return
	case util.Digit0:
		nextCH := s.Peek(1)
		// '0x', '0X' - 16进制
		if nextCH == util.LowercaseX || nextCH == util.UppercaseX {
			s.consumeRadixNumber(16)
			return
		}
		// '0o', '0O' - 8进制
		if nextCH == util.LowercaseO || nextCH == util.UppercaseO {
			s.consumeRadixNumber(8)
			return
		}
		// '0b', '0B' - 2进制
		if nextCH == util.LowercaseB || nextCH == util.UppercaseB {
			s.consumeRadixNumber(2)
			return
		}
		return
	case
		util.Digit1,
		util.Digit2,
		util.Digit3,
		util.Digit4,
		util.Digit5,
		util.Digit6,
		util.Digit7,
		util.Digit8,
		util.Digit9:
		s.consumeNumber(false)
		return
	case util.QuestionMark, util.Apostrophe:
		s.consumeQuotationMarkOrApostrophe()
		return
	case util.Slash:
		s.consumeSlash()
		return
	case util.PercentSign, util.Asterisk:
		s.consumePercentSignOrAsterisk()
		return
	case util.VerticalBar, util.Ampersand:
		s.consumeVerticalBarOrAmpersand()
		return
	case util.Caret:
		s.consumeCaret()
		return
	case util.PlusSign, util.Dash:
		s.consumePlusSignOrDash()
		return
	case util.LessThan:
		s.consumeLessThan()
		return
	case util.GreaterThan:
		s.consumeGreaterThan()
		return
	case util.EqualsTo, util.ExclamationMark:
		s.consumeExclamationMarkOrEqualsTo()
		return
	case util.Tilde:
		s.consumeTilde()
		return
	case util.AtSign:
		s.consumeAtSign()
		return
	case util.NumberSign:
		s.consumeNumberSign()
		return
	case util.Backslash:
		s.consumeDot()
		return
	default:
		//TODO: return error

	}
	return
}

// Peek 查看当前byte
func (s Lexer) Peek(n int) byte {
	return s.payload[s.index+n]
}

// Move 移动解析的位置
func (s Lexer) Move(n int) byte {
	if pos := n + s.index; s.len > pos && pos >= 0 {
		s.index = pos
	}
	return s.payload[s.index]
}

// Line 获得当前的行号
func (s Lexer) Line() (int, int) {
	return s.startLine, s.endLine
}

// Column 获取当前的列
func (s Lexer) Column() (int, int) {
	return s.startColumn, s.endColumn
}
func (s Lexer) Location() (int, int, int, int) {
	return s.startLine, s.endLine, s.startColumn, s.endColumn
}

// 跳过空格和换行
func (s Lexer) skipSpaceAndLineTerminator() {
	for {
		ch := s.payload[s.index]
		if isSpace(rune(ch)) {
			s.index++
			s.startColumn++
			s.endColumn++
			continue
		}
		if isLineTerminator(rune(ch)) {
			s.index++
			s.startLine++
			s.startColumn = 0
			s.endColumn = 0
			continue
		}
		return
	}
}

//消费.
func (s Lexer) consumeDot() {
	nextCH := s.Peek(1)
	//处理小数
	if nextCH >= util.Digit0 && nextCH <= util.Digit9 {
		s.consumeNumber(true)
		return
	}
	//处理扩展运算符
	if nextCH == util.Dot && s.Peek(2) == util.Dot {
		s.finishToken(types.EllipsisToken, 3)
		return
	}
	//处理正常的.
	s.finishToken(types.DoToken, 1)
	return
}

//消费number
func (s Lexer) consumeNumber(startsWithDot bool) {

}

// 消费string
func (s Lexer) consumeString() {

}

// 消费/
func (s Lexer) consumeSlash() {

}

// 消费*或者%
func (s Lexer) consumePercentSignOrAsterisk() {

}

//消费 ｜ &
func (s Lexer) consumeVerticalBarOrAmpersand() {

}

// 消费^
func (s Lexer) consumeCaret() {

}

//消费+ -
func (s Lexer) consumePlusSignOrDash() {

}

//消费 <
func (s Lexer) consumeLessThan() {

}

//消费>
func (s Lexer) consumeGreaterThan() {

}

//消费！=
func (s Lexer) consumeExclamationMarkOrEqualsTo() {

}

//消费 ～
func (s Lexer) consumeTilde() {

}

//消费 @
func (s Lexer) consumeAtSign() {

}

//消费#
func (s Lexer) consumeNumberSign() {

}

//消费\
func (s Lexer) consumeBackSlash() {

}

//消费 模版字符串
func (s Lexer) consumeTemplate() {

}

//消费 "
func (s Lexer) consumeQuotationMark() {

}

// 消费 \
func (s Lexer) consumeBackslash() {

}

// 当前token完结
func (s Lexer) finishToken(tokenType types.TokenType, cost int) {
	s.value = s.payload[s.index : s.index+cost]
	s.index += cost
	s.endColumn += cost
	s.tokenType = tokenType
}

//消费带进制的数
func (s Lexer) consumeRadixNumber(radix int) {
	var isBigInt = false
	val := s.consumeInt(radix)
	if val == nil {
		//TODO：错误处理
	}
	if isBigInt {

	}
}

//消费int
func (s Lexer) consumeInt(radix int) []byte {
	pos := s.index
	//如果为进制数需要跳过两位
	if radix != 10 {
		pos++
	}

	return nil
}

//消费 " '
func (s Lexer) consumeQuotationMarkOrApostrophe() {

}

// 报错
func (s Lexer) throwErr() {

}
