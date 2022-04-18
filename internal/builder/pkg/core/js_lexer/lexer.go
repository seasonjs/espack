package js_lexer

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/seasonjs/espack/internal/builder/pkg/core/types"
	"github.com/seasonjs/espack/internal/builder/pkg/core/util"
	"github.com/seasonjs/espack/internal/logger"
	"io"
	"io/ioutil"
	"strconv"
)

// 禁止的数字边界
type sForbiddenNumericSeparatorSiblings struct {
	decBinOct []byte
	hex       []byte
}

// 允许的数字边界
type sAllowedNumericSeparatorSiblings struct {
	bin []byte
	oct []byte
	dec []byte
	hex []byte
}

//没有实际意义只是绑定接口

var (
	forbiddenNumericSeparatorSiblings = sForbiddenNumericSeparatorSiblings{
		decBinOct: []byte{
			util.Dot,
			util.UppercaseB,
			util.UppercaseE,
			util.UppercaseO,
			util.Underscore, // multiple separators are not allowed
			util.LowercaseB,
			util.LowercaseE,
			util.LowercaseO,
		},
		hex: []byte{
			util.Dot,
			util.UppercaseX,
			util.Underscore, // multiple separators are not allowed
			util.LowercaseX,
		},
	}
	allowedNumericSeparatorSiblings = sAllowedNumericSeparatorSiblings{
		bin: []byte{
			util.Digit0,
			util.Digit1,
		},
		oct: []byte{
			util.Digit0,
			util.Digit1,
			util.Digit2,
			util.Digit3,
			util.Digit4,
			util.Digit5,
			util.Digit6,
			util.Digit7,
		},
		dec: []byte{
			util.Digit0,
			util.Digit1,
			util.Digit2,
			util.Digit3,
			util.Digit4,
			util.Digit5,
			util.Digit6,
			util.Digit7,
			util.Digit8,
			util.Digit9,
		},
		hex: []byte{
			util.Digit0,
			util.Digit1,
			util.Digit2,
			util.Digit3,
			util.Digit4,
			util.Digit5,
			util.Digit6,
			util.Digit7,
			util.Digit8,
			util.Digit9,
			util.UppercaseA,
			util.UppercaseB,
			util.UppercaseC,
			util.UppercaseD,
			util.UppercaseE,
			util.UppercaseF,
			util.LowercaseA,
			util.LowercaseB,
			util.LowercaseC,
			util.LowercaseD,
			util.LowercaseE,
			util.LowercaseF,
		},
	}
)

// Lexer is a lexer for JavaScript.
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
	err         []error         //当前遇到的错误
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

// https://tc39.github.io/ecma262/#sec-line-terminators
func isNewLine(code rune) bool {
	switch code {
	case rune(util.LineFeed), rune(util.CarriageReturn), util.LineSeparator, util.ParagraphSeparator:
		return true

	default:
		return false
	}
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
		//TODO 考虑 [|
		s.finishToken(types.OpenBracketToken, 1)
		return
	case util.RightSquareBracket:
		s.finishToken(types.CloseBracketToken, 1)
		return
	case util.LeftCurlyBrace:
		//TODO 考虑 {|
		s.finishToken(types.OpenBraceToken, 1)
		return
	case util.RightCurlyBrace:
		s.finishToken(types.CloseBraceToken, 1)
		return
	case util.Colon:
		//TODO 考虑function bind处理
		s.finishToken(types.ColonToken, 1)
		return
	case util.QuestionMark:
		s.consumeQuestionMark()
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
	case util.QuotationMark, util.Apostrophe:
		s.consumeString(rune(ch))
		return
	case util.Slash:
		s.consumeSlash()
		return
	case util.PercentSign, util.Asterisk:
		s.consumePercentSignOrAsterisk(ch)
		return
	case util.VerticalBar, util.Ampersand:
		s.consumeVerticalBarOrAmpersand(ch)
		return
	case util.Caret:
		s.consumeCaret()
		return
	case util.PlusSign, util.Dash:
		s.consumePlusSignOrDash(ch)
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
	start := s.index
	isFloat := false
	isBigInt := false
	isDecimal := false
	hasExponent := false
	isOctal := false
	if i, _ := s.consumeInt(10); startsWithDot && i == nil {
		//TODO error is not float
		return
	}
	hasLeadingZero :=
		s.index-start >= 2 &&
			s.Peek(start) == util.Digit0

	if hasLeadingZero {
		integer := s.payload[start : s.index-start]
		isOctal = hasLeadingZero && !bytes.ContainsAny(integer, "89")
	}
	nextCh := s.Peek(1)
	if nextCh == util.Dot && !isOctal {
		s.index = s.index + 1
		s.consumeInt(10)
		isFloat = true
		nextCh = s.Peek(0)
	}

	if (nextCh == util.UppercaseE || nextCh == util.LowercaseE) && !isOctal {
		s.index++
		nextCh = s.Peek(0)
		if nextCh == util.PlusSign || nextCh == util.Dash {
			s.index++
		}
		if i, _ := s.consumeInt(10); i == nil {
			//TODO error is not float
			return
		}
		isFloat = true
		hasExponent = true
		nextCh = s.Peek(0)
	}
	if nextCh == util.LowercaseN {
		if isFloat || hasLeadingZero {
			//TODO error big int
		}
		s.index++
		isBigInt = true
	}

	if nextCh == util.LowercaseM {
		if hasExponent || hasLeadingZero {

			//TODO error hasExponent or hasLeadingZero
		}
		s.index++
		isDecimal = true
	}

	//TODO 处理 IsIdentifierStart

	str := s.payload[start : s.index-start]
	str = bytes.Trim(str, "_mn")

	if isBigInt {
		s.finishToken(types.BigIntToken, len(str))
		return
	}

	if isDecimal {
		s.finishToken(types.DecimalToken, len(str))
		return
	}

	s.finishToken(types.NumericToken, len(str))
}

// 消费string
func (s Lexer) consumeString(quote rune) {
	var out []byte
	s.index++
	strStart := s.index
	for {
		if s.index >= s.len {
			//TODO: string 溢出
		}
		ch := s.Peek(0)
		if rune(ch) == quote {
			break
		}
		if ch == util.Backslash {
			//out += this.input.slice(chunkStart, this.state.pos);
			copy(out, s.payload[strStart:s.index])
			// $FlowFixMe
			copy(out, s.consumeEscapedChar(false))
			strStart = s.index
		} else if rune(ch) == util.LineSeparator || rune(ch) == util.ParagraphSeparator {
			s.index++
			s.startLine++
			s.startColumn = s.index
		} else if isNewLine(rune(ch)) {
			//throw this.raise(Errors.UnterminatedString, {
			//at: this.state.startLoc,
			//});

		} else {
			s.index++
		}
	}
	s.index++
	copy(out, s.payload[strStart:s.index])
	s.finishToken(types.StringToken, len(out))
}
func (s Lexer) consumeEscapedChar(inTemplate bool) []byte {
	return nil
}

// 消费/
func (s Lexer) consumeSlash() {
	nextCh := s.Peek(1)
	if nextCh == util.EqualsTo {
		s.finishToken(types.SlashAssignToken, 2)
	} else {
		s.finishToken(types.SlashToken, 1)
	}
}

// 消费*或者%
func (s Lexer) consumePercentSignOrAsterisk(code byte) {
	var tt types.TokenType
	size := 1
	nextCh := s.Peek(1)
	if code == util.Asterisk {
		tt = types.AsteriskToken
	} else {
		tt = types.PercentSignToken
	}
	// '**'
	if code == util.Asterisk && nextCh == util.Asterisk {
		size++
		nextCh = s.Peek(2)
		tt = types.ExponentToken
	}

	// '%=' , '*='
	if nextCh == util.EqualsTo {
		size++
		if code == util.PercentSign {
			tt = types.ModEqToken
		} else {
			tt = types.MulEqToken
		}
	}
	s.finishToken(tt, size)
}

//消费 ｜ &
func (s Lexer) consumeVerticalBarOrAmpersand(code byte) {
	// '||' '&&' '||=' '&&='
	var tt types.TokenType
	nextCh := s.Peek(1)
	if nextCh == code {
		if s.Peek(2) == util.EqualsTo {
			//||=
			if code == util.VerticalBar {
				s.finishToken(types.OrEqToken, 3)
			}
			if code == util.Ampersand {
				s.finishToken(types.AndEqToken, 3)
			}
		} else {
			if code == util.VerticalBar {
				tt = types.LogicalORToken
			} else {
				tt = types.LogicalANDToken
			}
			s.finishToken(tt, 2)
		}
		return
	}

	if code == util.VerticalBar {
		// '|>'
		if nextCh == util.GreaterThan {
			s.finishToken(types.PipelineToken, 2)
			return
		}
	}

	if nextCh == util.EqualsTo {
		s.finishToken(types.BitOrEqToken, 2)
		return
	}
	if code == util.VerticalBar {
		tt = types.BitOrToken
	} else {
		tt = types.BitAndToken
	}
	s.finishToken(tt, 1)

}

// 消费^
func (s Lexer) consumeCaret() {
	nextCh := s.Peek(1)
	// '^='
	if nextCh == util.EqualsTo {
		s.finishToken(types.BitXorAssignToken, 2)
	} else {
		s.finishToken(types.BitXorToken, 1)
	}
}

//消费+ -
func (s Lexer) consumePlusSignOrDash(code byte) {
	// '+-'
	nextCh := s.Peek(1)
	var tt types.TokenType
	if nextCh == code {
		if code == util.PlusSign {
			tt = types.IncrToken
		}
		if code == util.Dash {
			tt = types.DecrToken
		}
		s.finishToken(tt, 2)
		return
	}

	if nextCh == util.EqualsTo {
		if code == util.PlusSign {
			tt = types.AddAssignToken
		}
		if code == util.Dash {
			tt = types.SubAssignToken
		}
		s.finishToken(tt, 2)
	} else {
		if code == util.PlusSign {
			tt = types.AddToken
		}
		if code == util.Dash {
			tt = types.SubToken
		}
		s.finishToken(tt, 1)
	}
}

//消费 <
func (s Lexer) consumeLessThan() {
	// '<'
	nextCh := s.Peek(1)
	if nextCh == util.LessThan {
		//<<=
		if s.Peek(2) == util.EqualsTo {
			s.finishToken(types.LtLtAssignToken, 3)
			return
		}
		//<<
		s.finishToken(types.LtLtToken, 2)
		return
	}

	if nextCh == util.EqualsTo {
		// <=
		s.finishToken(types.LtAssignToken, 2)
		return
	}

	s.finishToken(types.LtToken, 1)
}

//消费>
func (s Lexer) consumeGreaterThan() {
	if s.Peek(1) == util.GreaterThan {
		if s.Peek(2) == util.GreaterThan {
			if s.Peek(3) == util.EqualsTo {
				//>>>=
				s.finishToken(types.GtGtGtEqToken, 4)
			}
			//>>>
			s.finishToken(types.GtGtGtToken, 3)
		} else {
			//>>=
			if s.Peek(2) == util.EqualsTo {
				s.finishToken(types.GtGtEqToken, 3)
			}
			//>>
			s.finishToken(types.GtGtToken, 2)
		}
	} else {
		//>
		s.finishToken(types.GtToken, 1)
	}
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

//消费 ?
func (s Lexer) consumeQuestionMark() {
	nextCH := s.Peek(1)
	nextTCH := s.Peek(2)
	//??
	if nextCH == util.QuestionMark {
		//??=
		if nextTCH == util.EqualsTo {
			s.finishToken(types.NullishEqToken, 3)
		} else {
			s.finishToken(types.NullishToken, 2)
		}
	} else if nextCH == util.Dot &&
		!(nextTCH >= util.Digit0 && nextTCH <= util.Digit9) {
		//?.
		s.finishToken(types.OptChainToken, 2)
	} else {
		s.finishToken(types.QuestionToken, 1)
	}
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
	isBigInt := false
	s.index = s.index + 2
	val, cost := s.consumeInt(radix)
	if val == nil {
		//TODO：错误处理
		panic("err")
	}
	if s.Peek(1) == util.LowercaseO {
		s.index++
		isBigInt = true
	} else if s.Peek(1) == util.LowercaseM {
		//TODO：错误处理
		panic("err")
	}
	if isBigInt {
		//TODO 处理jsBigInt
	}
	s.finishToken(types.NumericToken, cost)
}

//消费int radix = 2,8,10,16
func (s Lexer) consumeInt(radix int) ([]byte, int) {
	startPos := s.index
	//如果为进制数需要跳过两位
	var forbiddenSiblings, allowedSiblings []byte
	//确定上下边界
	switch radix {
	//二进制
	case 2:
		forbiddenSiblings = forbiddenNumericSeparatorSiblings.decBinOct
		allowedSiblings = allowedNumericSeparatorSiblings.bin
		//八进制
	case 8:
		forbiddenSiblings = forbiddenNumericSeparatorSiblings.decBinOct
		allowedSiblings = allowedNumericSeparatorSiblings.oct
	//16进制
	case 10:
		forbiddenSiblings = forbiddenNumericSeparatorSiblings.decBinOct
		allowedSiblings = allowedNumericSeparatorSiblings.dec
	case 16:
		forbiddenSiblings = forbiddenNumericSeparatorSiblings.hex
		allowedSiblings = allowedNumericSeparatorSiblings.hex
	}
	//TODO: js的number比int大需要进行转化
	total := 0
	//仅仅记录循环，而不打破
	for {
		ch := s.Peek(0)
		var val int

		if ch == util.Underscore {
			prev := s.Peek(-1)
			next := s.Peek(1)
			if !bytes.ContainsRune(allowedSiblings, rune(next)) {
				//TODO 报错
				panic("error:")
			} else if bytes.ContainsRune(forbiddenSiblings, rune(prev)) || bytes.ContainsRune(forbiddenSiblings, rune(next)) {
				//TODO 报错
				panic("error:")
			}
			//TODO 不允许_
			s.index++
		}
		if ch >= util.LowercaseA {
			val = int(ch - util.LowercaseA + util.LineFeed)
		} else if ch >= util.UppercaseA {
			val = int(ch - util.UppercaseA + util.LineFeed)
		} else if util.IsDigit(ch) {
			val = int(ch - util.Digit0) // 0-9
		}
		////TODO: 处理无穷
		//else {
		//
		//}
		if val >= radix {
			//TODO： 出错，溢出边界
			break
		}
		s.index++
		total = total*radix + val
		if s.index == startPos {
			return nil, 0
		}
	}
	//return int to []byte
	return []byte(strconv.Itoa(total)), s.index - startPos
}

// 报错
func (s Lexer) throwErr(message string) {
	s.err = append(s.err, errors.New(message))
}
