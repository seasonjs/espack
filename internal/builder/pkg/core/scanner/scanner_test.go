package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	//split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	//	//通过 bufio.ScanWords 这个扫描方法并不可行，需要根据需要定制开发
	//	return ScanKeyWords(data, atEOF)
	//}
	//sc.Split(split)
	line := 0
	for sc.Scan() {
		line++
		t.Log(sc.Text(), line)
	}

	if err := sc.Err(); err != nil {
		fmt.Printf("Invalid input: %s", err)
	}

}

func TestRead(t *testing.T) {
	//这样写就代表建立了一个大小为2的缓存区
	buff := make([]byte, 4)
	str := strings.NewReader(`
		// this is an Comment 
		`)
	for {
		_, err := str.Read(buff)
		if err != nil {
			return
		}
		t.Log(string(buff))
		//end += n
		//loop++
	}
}

func TestScanKeyWords(t *testing.T) {
	//str := strings.NewReader(`
	//	// this is an Comment
	//	function a(){}
	//	`)
	//sc := NewScanner(str)
	//for sc.Scan() {
	//	if sc.Err() != nil {
	//		t.Log(sc.Err().Error())
	//		return
	//	}
	//	t.Log(sc.JsTokenLocation())
	//}
}
