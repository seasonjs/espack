package parser

import (
	"github.com/seasonjs/espack/internal/builder/pkg/core/in"
	"github.com/seasonjs/espack/internal/builder/pkg/core/lexer"
	"testing"
)

func TestES5(t *testing.T) {
	l := lexer.NewLexer(in.NewInputString("/*test */ \r\n import React from 'react';Whitespace import React from '../React';\r\n  if (React}) { console?.log('React has been import'); }"))
	t.Log(l)
}
func TestProgram_ParseProgram(t *testing.T) {
	p := NewProgram()
	p.ParseProgram(in.NewInputString("/*test */ \r\n import React from 'react';Whitespace import React from '../React';\r\n  if (React}) { console?.log('React has been import'); }"))
}
