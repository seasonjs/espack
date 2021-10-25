package in

import (
	"testing"
)

func TestSlice(t *testing.T) {
	//numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	//a := numbers[1:3:3]
	//t.Log(a,numbers)
	//s := "abc"
	//var a []byte
	//copy(a[:], s)
	i := NewInputString("asdadsa")
	for {
		i.Move(1)
		buf := i.Shift()
		t.Log(buf)
	}

}
