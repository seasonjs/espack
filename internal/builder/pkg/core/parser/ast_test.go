package parser

import (
	"encoding/json"
	"testing"
)

func TestJsonBuild(t *testing.T) {
	s := PrivateIdentifier{
		Name: "abc",
	}
	j, _ := json.Marshal(s)
	t.Log(j)
	t.Log("abc")
	t.Log(string(j))
}
