package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestFs_ConvertPath(t *testing.T) {
	str, _ := os.Getwd()
	fmt.Println("The path is ", str)
}
func TestEnv_Dev(t *testing.T) {
	////filename:=GetCurrentPath()
	//fmt.Println("The path is ", filename)
}
