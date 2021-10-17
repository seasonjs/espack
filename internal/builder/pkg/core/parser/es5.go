package parser

import (
	"fmt"
	"github.com/seasonjs/espack/internal/builder/pkg/core/in"
	"github.com/seasonjs/espack/internal/builder/pkg/core/lexer"
	"github.com/seasonjs/espack/internal/logger"
	"io"
)

//类型声明基于 1. https://github.com/estree/estree/blob/master/es5.md
//           2. https://github.com/cst/cst
//第二参考文档 https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Expressions_and_Operators
// babel 是在node 节点扩展属性，但是这个方法go是肯定不合理的，所以我们只能挨个解析
// 与babel 类似的流程，可以先解析成类型再递进的深入逐渐解析到最终类型
// babel 转义流程
// parseTopLevel->parseProgram-> parseBlockBody->parseBlockOrModuleBlockBody->loop(parseStatement)
//							|
//							->parseInterpreterDirective
//parseStatement-> parseDecorators
//				|
//				->parseStatementContent->switch(parseAnyStatement)
// esbuild
//Parse->newParser->NewLexer->toAST
// espack 暂定的流程
// Parse->NewLexer-> ParseProgram->ParseStatement

type JsType string

//这个地方的类型生命是否应该采用和tokenType一样的类型声明方式
const (
	CommentBlockType          JsType = "CommentBlock"
	CommentLineType           JsType = "CommentLine"
	IdentifierType            JsType = "Identifier"
	LiteralType               JsType = "Literal"
	ProgramType               JsType = "Program"
	ExpressionStatementType   JsType = "ExpressionStatement"
	BlockStatementType        JsType = "BlockStatement"
	EmptyStatementType        JsType = "EmptyStatement"
	DebuggerStatementType     JsType = "DebuggerStatement"
	WithStatementType         JsType = "WithStatement"
	ReturnStatementType       JsType = "ReturnStatement"
	LabeledStatementType      JsType = "LabeledStatement"
	BreakStatementType        JsType = "BreakStatement"
	ContinueStatementType     JsType = "ContinueStatement"
	IfStatementType           JsType = "IfStatement"
	SwitchStatementType       JsType = "SwitchStatement"
	SwitchCaseType            JsType = "SwitchCase"
	ThrowStatementType        JsType = "ThrowStatement"
	TryStatementType          JsType = "TryStatement"
	CatchClauseType           JsType = "CatchClause"
	WhileStatementType        JsType = "WhileStatement"
	DoWhileStatementType      JsType = "DoWhileStatement"
	ForStatementType          JsType = "ForStatement"
	ForInStatementType        JsType = "ForInStatement"
	FunctionDeclarationType   JsType = "FunctionDeclaration"
	VariableDeclarationType   JsType = "VariableDeclaration"
	VariableDeclaratorType    JsType = "VariableDeclarator"
	ThisExpressionType        JsType = "ThisExpression"
	ArrayExpressionType       JsType = "ArrayExpression"
	ObjectExpressionType      JsType = "ObjectExpression"
	PropertyType              JsType = "Property"
	FunctionExpressionType    JsType = "FunctionExpression"
	UnaryExpressionType       JsType = "UnaryExpression"
	UpdateExpressionType      JsType = "UpdateExpression"
	BinaryExpressionType      JsType = "BinaryExpression"
	AssignmentExpressionType  JsType = "AssignmentExpression"
	LogicalExpressionType     JsType = "LogicalExpression"
	MemberExpressionType      JsType = "MemberExpression"
	ConditionalExpressionType JsType = "ConditionalExpression"
	CallExpressionType        JsType = "CallExpression"
	NewExpressionType         JsType = "NewExpression"
	SequenceExpressionType    JsType = "SequenceExpression"

	//ES2015

	SuperType                    JsType = "super"
	SpreadElementType            JsType = "SpreadElement"
	ArrowFunctionExpressionType  JsType = "ArrowFunctionExpression"
	YieldExpressionType          JsType = "YieldExpression"
	ExpressionType               JsType = "Expression"
	TaggedTemplateExpressionType JsType = "TaggedTemplateExpression"
	TemplateElementType          JsType = "TemplateElement"
	ObjectPatternType            JsType = "ObjectPattern"
	ArrayPatternType             JsType = "ArrayPattern"
	RestElementType              JsType = "RestElement"
	AssignmentPatternType        JsType = "AssignmentPattern"
	ClassBodyType                JsType = "ClassBody"
	MethodDefinitionType         JsType = "MethodDefinition"
	ClassDeclarationType         JsType = "ClassDeclaration"
	ClassExpressionType          JsType = "ClassExpression"
	MetaPropertyType             JsType = "MetaProperty"
	ImportDeclarationType        JsType = "ImportDeclaration"
	ImportSpecifierType          JsType = "ImportSpecifier"
	ImportDefaultSpecifierType   JsType = "ImportDefaultSpecifier"
	ImportNamespaceSpecifierType JsType = "ImportNamespaceSpecifier"
	ExportNamedDeclarationType   JsType = "ExportNamedDeclaration"
	ExportSpecifierType          JsType = "ExportSpecifier"
	ExportDefaultDeclarationType JsType = "ExportDefaultDeclaration"
	ExportAllDeclarationType     JsType = "ExportAllDeclaration"

	//ES2017

	AwaitExpressionType JsType = "AwaitExpression"

	//ES2020

	ChainExpressionType  JsType = "ChainExpression"
	ImportExpressionType JsType = "ImportExpression"

	//ES2022

	PropertyDefinitionType JsType = "PropertyDefinition"
	PrivateIdentifierType  JsType = "PrivateIdentifier"
	StaticBlockType        JsType = "StaticBlock"
)

type KindType string

const VarKind KindType = "var"

// Comment 注释基本类型 根据 estree的已有讨论，https://github.com/estree/estree/issues/41
// 注释需要单独处理并挂载在Node上，这里的Comment 结构我会跟babel的保持基本一致
type Comment struct {
}

//===================================================================================

//Node objects
type Node struct {
	*lexer.Lexer
	jsT              JsType
	leadingComments  []Comment
	trailingComments []Comment
	innerComments    []Comment
}

// NewNode 所有节点程序都可以被视为Node
func NewNode(lex *lexer.Lexer) *Node {
	node := &Node{}
	node.Lexer = lex
	return node
}

// StartNode 从头开始识别为Node
func (n *Node) StartNode() *Node {
	node := &Node{}
	node.Lexer = n.Lexer
	node.Lexer.Loc = n.Loc
	return node
}

// StartNodeAt 从指定位置开始识别为Node 此方法会新建Node实例
func (n *Node) StartNodeAt(loc lexer.SourceLocation) *Node {
	node := &Node{}
	node.Lexer = n.Lexer
	node.Lexer.Loc = loc
	//TODO 对注释进行拷贝
	return node
}

// StartNodeAtNode StartLoc a new node with a previous node's location.
func (n *Node) StartNodeAtNode(node Node) *Node {
	return n.StartNodeAt(node.Loc)
}

func (n *Node) finishNodeAt(loc lexer.SourceLocation) {
	n.Loc = loc
}

//func (n *Node) finishNode() {
//	n.
//}

//==========================================================================================

// Expression 表达式
type Expression struct{ Node }

func StartExpression(node Node) *Expression {
	return &Expression{
		node,
	}
}

//===========================================================================================

// Pattern 参数
type Pattern struct{ Node }

func StartPattern(node Node) Pattern {
	return Pattern{
		node,
	}
}

type StatementLike interface{}

type Statement struct{ *Node }

func StartStatement(node *Node) *Statement {
	return &Statement{
		node,
	}
}

func (s *Statement) ParseStatement() StatementLike {
	switch s.Cache.TT {
	case lexer.BreakToken:
		return StartBreakStatement(s).ParseBreakStatement()
	case lexer.ContinueToken:
		return StartContinueStatement(s).ParseContinueStatement()
	default:
		return nil
	}
}

//==================================================================================

// ParseStatement 解析为具体的Statement 暂时不考虑修饰器
//func (s *Statement) ParseStatement() {
//	node := Node{}
//}

// Identifier 变量名称，函数名称等一系列名称的定义
type Identifier struct {
	//Expression
	//Pattern
	*Node
	name string
}

func StartIdentifier(node *Node) *Identifier {
	id := &Identifier{}
	id.Node = node
	id.jsT = IdentifierType
	return id
}

func (i *Identifier) ParseIdentifier() *Identifier {
	//真的扫描到了 IdentifierToken
	if i.Cache.TT == lexer.IdentifierToken {
		i.name = string(i.Cache.Text)
	}
	//则说明其实没有Identifier
	return nil
}

//==================================================================================

type Literal struct {
	Expression
	value string
}

type RegExpLiteral struct {
	Literal
	regex struct {
		pattern string
		flags   string
	}
}

//==================================================================================

// Program ast 的入口
type Program struct {
	*Node
	body []interface{} //body [ Directive | Statement ]

}

func NewProgram(r *in.Input) *Program {
	//TODO 需要处理顶级注释
	topLevelNode := &Node{}

	l := lexer.NewLexer(r)
	topLevelNode.Lexer = l
	//不在顶层存储降低空间占用
	//topLevelNode.tokenValue = r.Bytes()
	//Line >=1
	topLevelNode.Loc.StartLoc.Line = 1
	return &Program{
		Node: topLevelNode,
	}
}

//func StartProgram(node *Node) *Program {
//	node.jsT = ProgramType
//	return &Program{
//		Node: node,
//	}
//}

func (p *Program) ParseProgram() {
	// 这里因为没有调用文本解析所以初始化代码位置暂时不检查
	//topLevelNode.finishNodeAt(r.Len(), Position{1, 0})
	//p = StartProgram(&topLevelNode)

	//TODO 处理前置指令
	// 比如 #! node
	//tt, text, loc := l.Next()
	//if tt == lexer.CommaToken {
	//
	//}
	for {
		p.Next()
		if p.Cache.TT == lexer.ErrorToken {
			if p.Err() != io.EOF {
				logger.Fail(fmt.Errorf("%s:%s:%v", p.Cache.Text, p.Err(), p.Cache.Loc), "Error on line")
			}
			return
		}
		node := p.StartNode()
		stmt := StartStatement(node).ParseStatement()
		p.body = append(p.body, stmt)
	}

}

//==================================================================================

// Directive 存放指令代码类似 "use strict"这样的 或者是#！这样的就是指令
// TODO 增加对指令文件的判断
type Directive struct {
	expression Literal
	directive  string
}

type FunctionBody struct {
	BlockStatement
	body interface{} //body  [ Directive | Statement ]

}

type BlockStatement struct {
	Statement
	body []Statement
}

type Function struct {
	Node
	id     Identifier
	params []Pattern
	body   FunctionBody
}

func (s *Statement) ParseFunction() {

}

type DebuggerStatement struct {
	Statement
}

type WithStatement struct {
	Statement
	object Expression
	body   Statement
}

type ReturnStatement struct {
	Statement
	argument Expression
}

type LabeledStatement struct {
	Statement
	label Identifier
	body  Statement
}

type BreakStatement struct {
	*Statement
	label *Identifier
}

func StartBreakStatement(s *Statement) *BreakStatement {
	s.jsT = BreakStatementType
	return &BreakStatement{
		Statement: s,
	}
}

// ParseBreakStatement 解析 break 关键字
// break [label];
// label 可选
// 与语句标签相关联的标识符。如果 break 语句不在一个循环或 switch 语句中，则该项是必须的。
func (s *BreakStatement) ParseBreakStatement() BreakStatement {
	s.Next()
	// 如果没有; 这意味着Break有label
	if s.Cache.TT != lexer.LineTerminatorToken {
		n := s.StartNode()
		ider := StartIdentifier(n).ParseIdentifier()
		s.label = ider
	}
	return *s
}

type ContinueStatement struct {
	*Statement
	label Identifier
}

func StartContinueStatement(s *Statement) *ContinueStatement {
	s.jsT = ContinueStatementType
	return &ContinueStatement{
		Statement: s,
	}
}

// ParseContinueStatement 解析 continue 关键字;
// continue [ label ];
// label
// 标识标号关联的语句
func (c *ContinueStatement) ParseContinueStatement() ContinueStatement {
	// 这意味着Break有label
	if c.Cache.TT != lexer.LineTerminatorToken {
		n := c.StartNode()
		ider := StartIdentifier(n).ParseIdentifier()
		c.label = *ider
	}
	return *c
}

type IfStatement struct {
	Statement
	test       Expression
	consequent Statement
	alternate  Statement
}

type SwitchStatement struct {
	Statement
	discriminant Expression
	cases        []SwitchCase
}

type SwitchCase struct {
	Node
	test       Expression
	consequent []Statement
}

type ThrowStatement struct {
	Statement
	argument Expression
}

type TryStatement struct {
	Statement
	block     BlockStatement
	handler   CatchClause
	finalizer BlockStatement
}

type CatchClause struct {
	Node
	param Pattern
	body  BlockStatement
}

type WhileStatement struct {
	Statement
	test Expression
	body Statement
}

type DoWhileStatement struct {
	Statement
	body Statement
	test Expression
}

type ForStatement struct {
	Statement
	//init: VariableDeclaration | Expression;
	init   interface{}
	test   Expression
	update Expression
	body   Statement
}

type ForInStatement struct {
	Statement
	//left VariableDeclaration |  Pattern;
	left  interface{}
	right Expression
	body  Statement
}

type Declaration struct{ Statement }

type FunctionDeclaration struct {
	Function
	Declaration
	id Identifier
}

type VariableDeclaration struct {
	Declaration
	declarations []VariableDeclarator
	kind         KindType
}

type VariableDeclarator struct {
	Node
	id   Pattern
	init Expression
}

type ThisExpression struct{ Expression }

type ArrayExpression struct {
	Expression
	elements []Expression
}

type ObjectExpression struct {
	Expression
	properties []Property
}

type Property struct {
	Node
	//key: Literal | Identifier;
	key   interface{}
	value Expression
	kind  KindType
}

type FunctionExpression struct {
	Function
	Expression
}

type UnaryExpression struct {
	Expression
	operator UnaryOperator
	prefix   bool
	argument Expression
}

// UnaryOperator "-" | "+" | "!" | "~" | "typeof" | "void" | "delete"
type UnaryOperator string

const (
	UnaryNegationOperator UnaryOperator = "-"
	UnaryPlusOperator     UnaryOperator = "+"
	NotOperator           UnaryOperator = "!"
	BitwiseNOTOperator    UnaryOperator = "~"
	TypeofOperator        UnaryOperator = "typeof"
	VoidOperator          UnaryOperator = "void"
	DeleteOperator        UnaryOperator = "delete"
)

type UpdateExpression struct {
	Expression
	operator UpdateOperator
	argument Expression
	prefix   bool
}

// UpdateOperator "++" | "--"
type UpdateOperator string

const (
	IncrementOperator UpdateOperator = "++"
	DecrementOperator UpdateOperator = "--"
)

type BinaryExpression struct {
	Expression
	operator BinaryOperator
	left     Expression
	right    Expression
}

//enum BinaryOperator {
//"==" | "!=" | "===" | "!=="
//| "<" | "<=" | ">" | ">="
//| "<<" | ">>" | ">>>"
//| "+" | "-" | "*" | "/" | "%"
//| "|" | "^" | "&" | "in"
//| "instanceof"
//}
type BinaryOperator string

const (
	EqualOperator          BinaryOperator = "=="
	NotEqualOperator       BinaryOperator = "!="
	StrictEqualOperator    BinaryOperator = "==="
	StrictNotEqualOperator BinaryOperator = "!=="

	LessThanOperator           BinaryOperator = "<"
	LessThanOrEqualOperator    BinaryOperator = "<="
	GreaterThanOperator        BinaryOperator = ">"
	GreaterThanOrEqualOperator BinaryOperator = ">="

	LeftShiftOperator                 BinaryOperator = "<<"
	SignPropagatingRightShiftOperator BinaryOperator = ">>"
	ZeroFillRightShift                BinaryOperator = ">>>"

	AdditionOperator       BinaryOperator = "+"
	SubtractionOperator    BinaryOperator = "-"
	MultiplicationOperator BinaryOperator = "*"
	DivisionOperator       BinaryOperator = "/"
	RemainderOperator      BinaryOperator = "%"

	BitwiseOROperator  BinaryOperator = "|"
	BitwiseXOROperator BinaryOperator = "^"
	BitwiseANDOperator BinaryOperator = "&"

	InOperator         BinaryOperator = "in"
	InstanceofOperator BinaryOperator = "instanceof"
)

type AssignmentExpression struct {
	Expression
	operator AssignmentOperators
	//left Pattern | Expression;
	left  interface{}
	right Expression
}

//enum AssignmentOperator {
//"=" | "+=" | "-=" | "*=" | "/=" | "%="
//| "<<=" | ">>=" | ">>>="
//| "|=" | "^=" | "&="
//}
type AssignmentOperators string

const (
	AssignmentOperator           AssignmentOperators = "="
	AdditionAssignment           AssignmentOperators = "+="
	SubtractionAssignment        AssignmentOperators = "-="
	MultiplicationAssignment     AssignmentOperators = "*="
	DivisionAssignment           AssignmentOperators = "/="
	RemainderAssignment          AssignmentOperators = "%="
	LeftShiftAssignment          AssignmentOperators = "<<="
	RightShiftAssignment         AssignmentOperators = ">>="
	UnsignedRightShiftAssignment AssignmentOperators = ">>>="
	BitwiseORAssignment          AssignmentOperators = "|="
	BitwiseXORAssignment         AssignmentOperators = "^="
	BitwiseANDAssignment         AssignmentOperators = "&="
)

type LogicalExpression struct {
	Expression
	operator LogicalOperator
	left     Expression
	right    Expression
}

//enum LogicalOperator {
//"||" | "&&"
//}

type LogicalOperator string

const (
	LogicalOROperator  LogicalOperator = "||"
	LogicalANDOperator LogicalOperator = "&&"
)

type MemberExpression struct {
	Expression
	Pattern
	object   Expression
	property Expression
	computed bool
}

type ConditionalExpression struct {
	Expression
	test       Expression
	alternate  Expression
	consequent Expression
}

type CallExpression struct {
	Expression
	callee    Expression
	arguments []Expression
}

type NewExpression struct {
	Expression
	callee    Expression
	arguments []Expression
}

type SequenceExpression struct {
	Expression
	expressions []Expression
}
