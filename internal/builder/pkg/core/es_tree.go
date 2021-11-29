// 这里主要参考了
// go parse 项目 https://github.com/tdewolff/parse/tree/master/js
//
// babel 原理文章 https://segmentfault.com/a/1190000017879365
// 以及 babel 项目
// esbuild 项目  https://github.com/evanw/esbuild
// estree https://github.com/estree/estree
// 使用estree是非常重要的，因为这意味着标准与一致性

// Lexer 暂时使用github.com/tdewolff/parse工具进行生成 但是如果要生成jsx 和ts 这个需要进行改造 ：）
// 更新一下，放弃直接使用 go parse 因为它缺少 token前移动 代码位置等信息 但是它的类型推断的使用是值得借鉴的
// 根据estree 进行分组继承可以有效实现babel相应功能并且会为兼容babel预留设计空间

package core

import (
	"github.com/seasonjs/espack/internal/builder/pkg/core/input"
	"github.com/seasonjs/espack/internal/builder/pkg/core/lexer"
	"github.com/seasonjs/espack/internal/builder/pkg/core/parser"
)

type ESVersion int

const (
	VES5 ESVersion = iota
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

type ESExtends uint

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
	Reader(path string, input *input.Input)
	Lexer()     // 词法分析器
	Parser()    //语法分析器
	Transform() //对语法进行转换
	Polyfills() //进行降级适配 core.js？
	Generator() //最终生成目标语法
}
type Tree struct {
	JSTarget ESVersion //需要支持的es版本
	Inputs   map[string]*input.Input
	Extends  []ESExtends
	ast      map[string]parser.Program
	lexers   map[string]*lexer.Lexer
}

func DefaultTree() ESTree {
	return Tree{
		JSTarget: VES2015, // 默认编译版本到es6
		Extends:  []ESExtends{JSON, JSX},
	}
}

func (t Tree) Reader(path string, input *input.Input) {
	panic("implement me")
}

func (t Tree) Lexer() {
	panic("implement me")
}

func (t Tree) Parser() {
	panic("implement me")
}

func (t Tree) Transform() {
	panic("implement me")
}

func (t Tree) Polyfills() {
	panic("implement me")
}

func (t Tree) Generator() {
	panic("implement me")
}
