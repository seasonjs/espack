package lexer_old

import (
	"fmt"
	"github.com/seasonjs/espack/internal/builder/pkg/core/input"
	"io"
	"testing"
)

func TestLexer(t *testing.T) {
	l := NewLexer(input_old.NewInputString(`
	function(a){
		console.log('test lexer')
	}
	`))
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
		default:
			t.Log("other", l.Cache.TT, string(l.Cache.Text), l.Cache.Loc)
		}
	}
}
