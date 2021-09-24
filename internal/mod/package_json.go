package mod

// TODO 需要理解package.json每个字段的概念 https://docs.npmjs.com/cli/v7/configuring-npm/package-json
type packageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// 每个包都要获取他的packagejson
func NewPackageJson() {

}
func (j *packageJSON) GetDependencies() {

}
