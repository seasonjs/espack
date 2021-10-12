package parser

//类型声明基于 https://github.com/estree/estree/blob/master/es5.md
//第二参考文档 https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Expressions_and_Operators
// babel 是在node 节点扩展属性，但是这个方法go是肯定不合理的，所以我们只能挨个解析
// 与babel 类似的流程，可以先解析成类型再递进的深入逐渐解析到最终类型

type JsType string

const (
	IdentifierType            JsType = "Identifier" // Identifier type
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
)

type KindType string

const VarKind KindType = "var"

type Position struct {
	line   int // >= 1
	column int // >= 0
}

type SourceLocation struct {
	//TODO 是否需要替换成[][]byte
	source string
	start  Position
	end    Position
}

//Node objects
type Node struct {
	jsT JsType
	loc SourceLocation
}

func (n *Node) name() {

}

type Expression struct{ Node }

type Pattern struct{ Node }

type Statement struct{ Node }

// ParseStatement 解析为具体的Statement 暂时不考虑修饰器
//func (s *Statement) ParseStatement() {
//	node := Node{}
//}

// Identifier 变量名称，函数名称等一系列名称的定义
type Identifier struct {
	Expression
	Pattern
	name string
}

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

// Program ast 的入口
type Program struct {
	Node
	//TODO 是否需要转化为map
	body []interface{} //body [ Directive | Statement ]

}

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
	Statement
	label Identifier
}

type ContinueStatement struct {
	Statement
	label Identifier
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
