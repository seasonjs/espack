package parser_old

import (
	"github.com/seasonjs/espack/internal/builder/pkg/core/input"
	"github.com/seasonjs/espack/internal/builder/pkg/core/lexer"
	"testing"
)

func TestES5(t *testing.T) {
	l := lexer.NewLexer(input.NewInputString("/*test */ \r\n import React from 'react';Whitespace import React from '../React';\r\n  if (React}) { console?.log('React has been import'); }"))
	t.Log(l)
}
func TestProgram_ParseProgram(t *testing.T) {
	p := NewProgram(input.NewInputString("/*test */ \r\n import React from 'react';Whitespace import React from '../React';\r\n  if (React}) { console?.log('React has been import'); }"))
	p.ParseProgram()
}
func TestBreak_Continue(t *testing.T) {
	//p := lexer.NewLexer(in.NewInputString("while (i < 6) {    if (i == 3) {      break;    }    i += 1;  }"))
	//for {
	//	p.Next()
	//	if p.Cache.TT == lexer.ErrorToken {
	//		if p.Err() != io.EOF {
	//			t.Log(p.Cache.Text, p.Err(), p.Cache.Loc)
	//			//logger.Fail(fmt.Errorf("%s:%s:%v", , "Error on line")
	//		}
	//		return
	//	}
	//	t.Log(p)
	//	//node := p.StartNode()
	//	//stmt := StartStatement(node).ParseStatement()
	//	//p.body = append(p.body, stmt)
	//}
	//
	p := NewProgram(input.NewInputString("while (i < 6) {\n    if (i == 3) {\n      break;\n    }\n    i += 1;\n  }"))
	p.ParseProgram()
	t.Log(p)

}
