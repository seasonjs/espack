package lexer

import (
	"fmt"
	"github.com/seasonjs/espack/internal/builder/pkg/core/in"
	"io"
	"testing"
)

func TestLexer(t *testing.T) {
	l := NewLexer(in.NewInputString("if (state == 5) { console.log('In state five'); }"))
	for {
		tt, text, loc := l.Next()
		switch tt {
		case ErrorToken:
			if l.Err() != io.EOF {
				fmt.Println("Error on line", text, ":", l.Err(), loc)
			}
			return
		case IdentifierToken:
			fmt.Println("Identifier", string(text), loc)
		case NumericToken:
			fmt.Println("Numeric", string(text), loc)
		default:
			t.Log("other", tt, string(text), loc)
		}
	}
}
