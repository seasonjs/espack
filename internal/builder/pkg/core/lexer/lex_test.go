package lexer

import (
	"fmt"
	"github.com/seasonjs/espack/internal/builder/pkg/core/in"
	"io"
	"testing"
)

func TestLexer(t *testing.T) {
	//l := NewLexer(in.NewInputString("/*test \r\n*/ \r\n import React from 'react';Whitespace import React from '../React';\r\n  if (React}) { console?.log('React has been import'); }"))
	//for {
	//	 l.Next()
	//	switch l.Cache.TT {
	//	case ErrorToken:
	//		if l.Err() != io.EOF {
	//			fmt.Println("Error on line", l.Cache.Text, ":", l.Err(), l.Cache.Loc)
	//		}
	//		return
	//	case IdentifierToken:
	//		fmt.Println("Identifier", string(l.Cache.Text), l.Cache.Loc)
	//	//case NumericToken:
	//	//	fmt.Println("Numeric", string(text), loc)
	//	default:
	//		t.Log("other", l.Cache.TT ,  string(l.Cache.Text), l.Cache.Loc)
	//	}
	//}
	l := NewLexer(in.NewInputString("outer_block:{\n\n  inner_block:{\n    console.log ('1');\n    break outer_block;      // breaks out of both inner_block and outer_block\n    console.log (':-(');    // skipped\n  }\n\n  console.log ('2');        // skipped\n}"))
	for {
		l.Next()
		switch l.Cache.TT {
		case ErrorToken:
			if l.Err() != io.EOF {
				fmt.Println("Error on line", l.Cache.Text, ":", l.Err(), l.Cache.Loc)
			}
			return
		case IdentifierToken:
			fmt.Println("Identifier", string(l.Cache.Text), l.Cache.Loc)
		//case NumericToken:
		//	fmt.Println("Numeric", string(text), loc)
		default:
			t.Log("other", l.Cache.TT, string(l.Cache.Text), l.Cache.Loc)
		}
	}
}
