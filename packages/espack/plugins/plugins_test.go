package plugins

//TODO 修改测试用例
//func TestPluginQueue_Add(t *testing.T) {
//	var que = NewPluginQueue()
//	que.Add(1)
//	if que.Len() != 1 {
//		t.Errorf("len expected be 1, but %d got\n", que.Len())
//	}
//	que.Add(2)
//	if que.Len() != 2 {
//		t.Errorf("len expected be 2, but %d got\n", que.Len())
//	}
//	if que.plugins[0] != 1 {
//		t.Errorf("que.plugins[0] expected be 1, but %d got\n", que.plugins[0])
//	}
//	if que.plugins[1] != 2 {
//		t.Errorf("que.plugins[0] expected be 1, but %d got\n", que.plugins[1])
//	}
//}
//func TestPluginQueue_Remove(t *testing.T) {
//	var que = NewPluginQueue()
//	v := que.Remove()
//	if v != nil {
//		t.Errorf("v expected be nil, but %d got\n", v)
//	}
//	que.Add(1)
//	que.Add(2)
//	//[1,2]
//	v = que.Remove()
//	//v=1 que = [2]
//	if v != 1 {
//		t.Errorf("v expected be 1, but %d got\n", v)
//	}
//	//v=2 que = []
//	v = que.Remove()
//	if v != 2 {
//		t.Errorf("v expected be 2, but %d got\n", v)
//	}
//}
//
//func TestPluginQueue_Next(t *testing.T) {
//	var que = NewPluginQueue()
//	que.Add(1)
//	que.Add(2)
//	//[1,2]
//	v := que.Next()
//	//1
//	if v != 1 {
//		t.Errorf("v expected be 1, but %d got\n", v)
//	}
//	v = que.Next()
//	//2
//	if v != 2 {
//		t.Errorf("v expected be 2, but %d got\n", v)
//	}
//	que.Remove()
//	//[2]
//	que.Remove()
//	//[]
//	v = que.Next()
//	if v != nil {
//		t.Errorf("v expected be nil, but %d got\n", v)
//	}
//	v = que.Next()
//	if que.cur > que.Len() && v != nil {
//		t.Errorf("v expected be nil, but %d got\n", v)
//	}
//	if que.cur < que.Len() && v == nil {
//		t.Errorf("v expected be not nil, but %d got\n", v)
//	}
//}
//func TestPluginQueue_Prev(t *testing.T) {
//	var que = NewPluginQueue()
//	v := que.Prev()
//	if v != nil {
//		t.Errorf("v expected be nil, but %d got\n", v)
//	}
//	que.Add(1)
//	que.Add(2)
//	//[1,2]
//	v = que.Prev()
//	//nil
//	if v != nil {
//		t.Errorf("v expected be nil, but %d got\n", v)
//	}
//	v = que.Prev()
//	//nil
//	if v != nil {
//		t.Errorf("v expected be nil, but %d got\n", v)
//	}
//	que.Remove()
//	//[2]
//	que.Next()
//	// 2
//	que.Next()
//	// nil
//	que.Prev()
//	//nil
//	if v != nil {
//		t.Errorf("v expected be nil, but %d got\n", v)
//	}
//	que.Add(3)
//	que.Add(4)
//	que.Next()
//	que.Next()
//	v = que.Prev()
//	if v != 3 {
//		t.Errorf("v expected be 2, but %d got\n", v)
//	}
//	v = que.Prev()
//	if v != 2 {
//		t.Errorf("v expected be 1, but %d got\n", v)
//	}
//}
