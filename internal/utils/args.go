package utils

import "time"

// 无实际意义，仅作为区分使用
type args struct {
}

// 校验参数是否可填写 args
func (a *args) CheckArgs(args ...interface{}) {
	time.Now()
}
