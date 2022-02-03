package js_parser

// 这里主要参考了
// go parse 项目 https://github.com/tdewolff/parse/tree/master/js
// babel 原理文章 https://segmentfault.com/a/1190000017879365
// 以及 babel 项目
// esbuild 项目  https://github.com/evanw/esbuild
// estree https://github.com/estree/estree
// 使用estree是非常重要的，因为这意味着标准与一致性

// Lexer 暂时使用github.com/tdewolff/parse工具进行生成 但是如果要生成jsx 和ts 这个需要进行改造 ：）
// 更新一下，放弃直接使用 go parse 因为它缺少 token前移动 代码位置等信息 但是它的类型推断的使用是值得借鉴的
// 根据estree 进行分组继承可以有效实现babel相应功能并且会为兼容babel预留设计空间

// 2022.1.13 更新一下，重新梳理了一下这里的逻辑，我觉得不应该是lexer操作移动然后在再生成语法树应该理解为：
// 需要有一个遍历到工具包buffer_scanner，这个buffer scanner 足够低级，只是为了能够对buffer实现类型string的操作
// 然后我们通过 scanner 进行逐行解析，最后我们在这个移动的过程中构造一个AST即可，也就是说我们会生成一个语法树
// 即我们都会生成ast，我们不会重复解析已经解析过到内容，但是我们可能会打断前进，返回到之前到位置，再次进行解析。
// @Note 在解析的过程中为了优化内存占用，使用buffer比较合理，但是这样还不如直接全部读取然后使用string，
//  生成的时候我们会频繁操作拼接所以也是使用buffer更合理

// 2022.2.3 之前想到到还是不完全，先将单个文件全都读到内存中，然后再进行解析，因为一个js文件再大也不可能过G？
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
