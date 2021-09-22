package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
)

// 无实际意义，仅作为区分使用
type err struct {
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// LogAndExit TODO: 待增加只报错不退出逻辑
func (e err) LogAndExit(msg error) {
	fmt.Println(msg.Error())
	if err, ok := msg.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			fmt.Printf("%+s:%d\n", f, f)
		}
	}
	os.Exit(1)
}
