package js_lexer

import (
	"github.com/seasonjs/espack/internal/builder/pkg/core/util"
	"github.com/seasonjs/espack/internal/logger"
	"io"
	"io/ioutil"
)

// 前置知识了解：https://go.dev/blog/strings
// TODO 需要保证payload不会溢出

type Lexer struct {
	payload  []byte //要扫描的代码
	index    int    //当前的位置
	len      int    // payload的总长度
	line     int    //当前代码行数
	column   int    //当前代码的列数
	ctxRange []int  //当前token的开始和结束位置
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
	s.ctxRange = []int{s.index, s.index + 1}
	ch := s.Peek(0)
	//判断是不是js的token
	switch ch {
	case util.Dot:
		s.consumeDot()
		return
	case util.LeftParenthesis:
		//TODO
		return
	case util.RightParenthesis:
		//TODO
		return
	case util.Semicolon:
		//TODO
		return
	case util.Comma:
		//TODO
		return
	case util.LeftSquareBracket:
		//TODO
		return
	case util.RightSquareBracket:
		//TODO
		return
	case util.LeftCurlyBrace:
		//TODO
		return
	case util.RightCurlyBrace:
		//TODO
		return
	case util.Colon:
		//TODO
		return
	case util.QuotationMark:
		s.consumeQuotationMark()
		return
	case util.GraveAccent:
		s.consumeTemplate()
		return
	case util.Digit0:
		//TODO
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
		s.consumeNumber()
		return
	case util.QuestionMark, util.Apostrophe:
		s.consumeQuotationMark()
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
func (s Lexer) Line() int {
	return s.line
}

// Column 获取当前的列
func (s Lexer) Column() int {
	return s.column
}

// 跳过空格和换行
func (s Lexer) skipSpaceAndLineTerminator() {
	for {
		ch := s.payload[s.index]
		if isSpace(rune(ch)) {
			s.index++
			s.column++
			continue
		}
		if isLineTerminator(rune(ch)) {
			s.index++
			s.line++
			s.column = 0
			continue
		}
		return
	}
}

//消费.
func (s Lexer) consumeDot() {

}

//消费number
func (s Lexer) consumeNumber() {

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
