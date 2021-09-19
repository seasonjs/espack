package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	cf := NewConfig()
	if *cf != (configuration{}) {
		t.Errorf(" cf expected be empty, but %v got\n", cf)
	}
}

func TestConfiguration_ReadFile(t *testing.T) {
	cf := NewConfig().ReadFile()
	if cf.Output != "pathxxxxxxxxxxxxx" {
		t.Errorf(" cf.Path expected be pathxxxxxxxxxxxxx, but %s got\n", cf.Output)
	}
	cf.ReadFile("./config.json")
	if cf.Entry != "ppppppp" {
		t.Errorf(" cf.Path expected be ppppppp, but %s got\n", cf.Entry)
	}
	//test error path panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	cf.ReadFile("./aaaaaa")

}
