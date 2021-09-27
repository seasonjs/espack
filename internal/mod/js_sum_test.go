package mod

import "testing"

func TestJsSum_FetchRootModMeta(t *testing.T) {
	sum := NewJsSum().FetchRootModMeta()
	t.Log(sum)
}
func TestJsSum_PrunerAndFetcher(t *testing.T) {
	sum := NewJsSum()
	sum.modList = append(sum.modList, &modInfo{})
	sum.modList = nil
	sum.modList = append(sum.modList, &modInfo{})
	t.Log(sum.modList)
}
func TestJsSum_PrunerAndFetcher2(t *testing.T) {
	js := NewJsSum().FetchRootModMeta()
	js.PrunerAndFetcher()
	t.Log(js)
}
