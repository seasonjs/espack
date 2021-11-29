package input_old

import (
	"io"
	"testing"
)

func TestSlice(t *testing.T) {
	i := NewInputString(`
	function a(){
		console.log('input ')
	}
	`)
	for {
		i.Move(1)
		buf := i.Shift()
		t.Log(string(buf))
		if i.Err() == io.EOF {
			return
		}
	}

}
