// 这里主要参考了
// go parse 项目 https://github.com/tdewolff/parse/tree/master/js
//
// babel 原理文章 https://segmentfault.com/a/1190000017879365
// esbuild 项目  https://github.com/evanw/esbuild
// estree https://github.com/estree/estree
// 使用estree是非常重要的，因为这意味着标准与一致性

// Lexer 暂时使用github.com/tdewolff/parse工具进行生成 但是如果要生成jsx 和ts 这个需要进行改造 ：）
// 更新一下，放弃直接使用 go parse 因为它缺少 token前移动 代码位置等信息 但是它的类型推断的使用是值得借鉴的
// 根据estree 进行分组继承可以有效实现babel相应功能并且会为兼容babel预留设计空间

package core

import (
	"fmt"
	"github.com/seasonjs/espack/internal/builder/pkg/core/in"
	"github.com/seasonjs/espack/internal/builder/pkg/core/lexer"
	"github.com/seasonjs/espack/internal/builder/pkg/core/parser_old"
	"github.com/seasonjs/espack/internal/logger"
)

type ESVersion int

const (
	VDefault ESVersion = iota
	VES5
	VES2015
	VES2016
	VES2017
	VES2018
	VES2019
	VES2020
	VES2021
	VES2022
	VTypeScript
)

type ESGVersion int

const (
	GDefault ESGVersion = iota
	GES5
	GES2015
	GES2016
	GES2017
	GES2018
	GES2019
	GES2020
	GES2021
	GES2022
)

type ESExtends int

const (
	JSX ESExtends = iota
	JSON
	//是否要通过estree支持以下
	//CSS in Js
	//SVG
	//PNG
	//VUE
)

type ESTree interface {
	Reader(path string, input *in.Input)
	Lexer()     // 词法分析器
	Parser()    //语法分析器
	Transform() //对语法进行转换
	Polyfills() //进行降级适配 core.js？
	Generator() //最终生成目标语法
}

type ES5 struct {
	inputs   map[string]*in.Input
	lexers   map[string]*lexer.Lexer
	gVersion ESGVersion
	extends  []ESExtends
	ast      map[string]parser_old.Program
}

func NewESTree(v ESVersion, gv ESGVersion, extends ...ESExtends) ESTree {
	versionToInstanceMap := map[ESVersion]ES5Interface{
		VES5:    &ES5{},
		VES2015: &ES2015{},
		VES2016: &ES2016{},
		VES2017: &ES2017{},
		VES2018: &ES2018{},
		VES2019: &ES2019{},
		VES2020: &ES2020{},
		VES2021: &ES2021{},
		VES2022: &ES2022{},
	}

	instance, ok := versionToInstanceMap[v]
	if ok {
		instance.SetTargetVersion(gv)
		instance.SetESExtends(extends...)
		return instance
	}

	// map里面没有找到 额外判断
	if v == VTypeScript {
		logger.Fail(fmt.Errorf("暂不支持对ts的解析"), "estree 错误终止解析")
		return nil
	}

	logger.Info("未定义解析版本,将启用默认设置和自动推断模式进行解析")
	//TODO 需要在这里进行推断
	return &ES2022{}
}

type ES5Interface interface {
	ESTree
	SetTargetVersion(ESGVersion)
	SetESExtends(...ESExtends)
}

//======================= getter setter ===========================

func (es *ES5) SetESExtends(v ...ESExtends) {
	es.extends = v
}

func (es ES5) GetESExtends() []ESExtends {
	return es.extends
}

func (es *ES5) SetTargetVersion(gv ESGVersion) {
	es.gVersion = gv
}

func (es ES5) GetTargetVersion() ESGVersion {
	return es.gVersion
}

func (es ES5) GetAST() map[string]parser_old.Program {
	return es.ast
}

func (es *ES5) SetAST(ast *map[string]parser_old.Program) {
	es.ast = *ast
}

//===================================================================

// 这意味着可以通过词法分析后 在转换为 ast 之前就逐步的将全部代码导入

func (es *ES5) Reader(path string, r *in.Input) {
	if es.inputs == nil {
		es.inputs = make(map[string]*in.Input)
	}
	es.inputs[path] = r
}

func (es *ES5) Lexer() {
	if es.lexers == nil {
		es.lexers = make(map[string]*lexer.Lexer)
	}
	for s, input := range es.inputs {
		es.lexers[s] = lexer.NewLexer(input)
	}
	//如果处理完后清理入口
	es.inputs = nil
}

func (es ES5) Parser() {
	//暂不考虑 process shebang
	//for s, lexer := range es.lexers {
	//	_, data := lexer.Next()
	//}

}

// Transform 在此处调用插件的 Transform
func (es *ES5) Transform() {

}

// Polyfills 在此处根据版本调用 Polyfills
func (es *ES5) Polyfills() {

}

// Generator 将ast重新转换成js
func (es *ES5) Generator() {

}

//=============================================================================

type ES2015 struct {
	ES5
	ast map[string]parser_old.ProgramES2015
}

func (es *ES2015) Parser() {

}

//======================= getter setter ========================================

func (es ES2015) GetAST() map[string]parser_old.ProgramES2015 {
	return es.ast
}

func (es *ES2015) SetAST(ast *map[string]parser_old.ProgramES2015) {
	es.ast = *ast
}

//=============================================================================

type ES2016 struct {
	ES2015
}

func (es *ES2016) Parser() {

}

//=============================================================================

type ES2017 struct {
	ES2016
}

func (es *ES2017) Parser() {

}

//=============================================================================

type ES2018 struct {
	ES2017
}

func (es *ES2018) Parser() {

}

//=============================================================================

type ES2019 struct {
	ES2018
}

func (es *ES2019) Parser() {

}

//=============================================================================

type ES2020 struct {
	ES2019
}

func (es *ES2020) Parser() {

}

//=============================================================================

type ES2021 struct {
	ES2020
}

func (es *ES2021) Parser() {

}

//=============================================================================

type ES2022 struct {
	ES2021
}

func (es *ES2022) Parser() {

}
