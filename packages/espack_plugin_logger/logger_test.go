package espack_plugin_logger

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func TestLogInfo(t *testing.T) {
	Info("%s", "11111")
	fmt.Println("????")
}
func TestTrace(t *testing.T) {
	err := errors.New("an error")
	Trace(err, "错误消息：")
}
