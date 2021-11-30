// Copyright (c) 2015 Taco de Wolff.
// Use of this source code is governed by a MIT style
// license that can be found in https://github.com/tdewolff/parse/blob/master/LICENSE.md

package scanner

import (
	"io"
	"unicode/utf8"
)

type Scanner struct {
	r      io.Reader
	buf    []byte //内部缓存的buf长度
	line   int    //行号
	column int    //列号
	pos    int    //byte 的index
	token  []byte //当前的token
	tt     string //当前的token的类型
	err    error  //遇到的错误
}

// NewScanner returns a new Scanner to read from r.
// The split function defaults to ScanLines.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r: r,
	}
}

// Err returns the first non-EOF error that was encountered by the Scanner.
func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

// Bytes bytes格式返回当前token
func (s *Scanner) Bytes() []byte {
	return s.token
}

// Text String 格式返回当前token
func (s *Scanner) Text() string {
	return string(s.token)
}

// isSpace 是否是空格
func isSpace(r rune) bool {
	if r <= '\u00FF' {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	// High-valued ones.
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}

// ScanKeyWords 分割代码关键词 如果是空格则去除
func ScanKeyWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSpace(r) {
			break
		}
	}
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if isSpace(r) {
			return i + width, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}
