package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestFile(t *testing.T) {
	js, err := filepath.Abs("../../../../case/index.jsx")
	if err != nil {
		return
	}
	open, err := os.Open(js)
	defer func(open *os.File) {
		err := open.Close()
		if err != nil {
			t.Log("close error")
		}
	}(open)
	sc := bufio.NewScanner(open)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		//通过 bufio.ScanWords 这个扫描方法并不可行，需要根据需要定制开发
		return ScanKeyWords(data, atEOF)
	}
	sc.Split(split)
	for sc.Scan() {
		fmt.Printf("%s\n", sc.Text())
	}

	if err := sc.Err(); err != nil {
		fmt.Printf("Invalid input: %s", err)
	}

}
