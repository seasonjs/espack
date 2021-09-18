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
	if cf.Path != "pathxxxxxxxxxxxxx" {
		t.Errorf(" cf.Path expected be pathxxxxxxxxxxxxx, but %s got\n", cf.Path)
	}
	cf.ReadFile("./config.json")
	if cf.Path != "ppppppp" {
		t.Errorf(" cf.Path expected be ppppppp, but %s got\n", cf.Path)
	}
	//test error path panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	cf.ReadFile("./aaaaaa")

}
