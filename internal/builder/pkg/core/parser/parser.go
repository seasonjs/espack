package parser

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

type Parser interface {
	Reader(files map[string][]byte) //资源读取
	//Lexer(lex *lexer_old.Lexer)               // 词法分析器

	//TODO: 增加之后的转换和填充

	//Transform() //对语法进行转换
	//Polyfills() //进行降级适配 core.js？
	//Generator() //最终生成目标语法
}
type DefaultParserImpl struct {
	JSTarget ESVersion //需要支持的es版本
	Inputs   map[string][]byte
	Extends  []ESExtends
	ast      map[string]Program
	//lexer    *lexer_old.Lexer
}

// DefaultParser 使用默认的解析器
func DefaultParser() Parser {
	return DefaultParserImpl{
		JSTarget: VES2015, // 默认编译版本到es6
		Extends:  []ESExtends{JSON, JSX},
	}
}

// Reader 设置要读取的文件
func (t DefaultParserImpl) Reader(files map[string][]byte) {
	t.Inputs = files
}

//// Lexer 设置Lexer
//func (t DefaultParserImpl) Lexer(lex *lexer_old.Lexer) {
//	t.lexer = lex
//}

//func (t DefaultParserImpl) Transform() {
//	panic("implement me")
//}
//
//func (t DefaultParserImpl) Polyfills() {
//	panic("implement me")
//}
//
//func (t DefaultParserImpl) Generator() {
//	panic("implement me")
//}
