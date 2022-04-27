package config

//func TestNewConfig(t *testing.T) {
//	cf := NewConfig()
//	if *cf != (Configuration{}) {
//		t.Errorf(" cf expected be empty, but %v got\n", cf)
//	}
//}
//
//func TestConfiguration_ReadFile(t *testing.T) {
//	cf := NewConfig()
//	path := utils.FS.GetCurrentPath()
//	path = filepath.Join(path, "../case/espack.config.json")
//	cf.ReadFile(path)
//	if cf.Entry != "ppppppp" {
//		t.Errorf(" cf.Path expected be ppppppp, but %s got\n", cf.Entry)
//	}
//	//test error path panic
//	//defer func() {
//	//	if r := recover(); r == nil {
//	//		t.Errorf("The code did not panic")
//	//	}
//	//}()
//	//cf.ReadFile("./aaaaaa")
//
//}
//func TestConfiguration_GetEntry(t *testing.T) {
//	cf := NewConfig()
//	path := utils.FS.GetCurrentPath()
//	path = filepath.Join(path, "../case/espack.config.json")
//	cf.ReadFile(path)
//	cf.FormatEntry().FormatOutput()
//	fmt.Println(cf.FormatEntry())
//}
