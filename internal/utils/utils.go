package utils

import "sync"

var (
	FS   *fs
	Args *args
	Env  *env
	Err  *err
	once sync.Once
)

func init() {
	once.Do(func() {
		Env = &env{}
		//newEnv()
		FS = &fs{}
		Args = &args{}
		Err = &err{}
	})
}
