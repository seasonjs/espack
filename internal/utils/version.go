package utils

type version struct {
}
type versionCtrl struct {
	version string
}

// CheckVersion 检查是否是合法版本
func (v version) CheckVersion(ver string) bool {
	return versionReg.MatchString(ver)
}

// FindVersionStr 找到版本字段
func (v version) FindVersionStr(ver string) map[string]string {
	match := versionMatchReg.FindStringSubmatch(ver)
	results := map[string]string{}
	for i, name := range match {
		results[versionReg.SubexpNames()[i]] = name
	}
	return results
}
