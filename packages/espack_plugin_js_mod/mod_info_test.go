package espack_plugin_js_mod

import "testing"

func TestModInfo_FetchModInfo(t *testing.T) {
	NewModInfo("react", "17.0.1").FetchModInfo()
}
