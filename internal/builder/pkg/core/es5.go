package core

//类型声明基于 https://github.com/estree/estree/blob/master/es5.md
//第二参考文档 https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Expressions_and_Operators

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

const VarType KindType = "var"

type Position struct {
	line   int // >= 1
	column int // >= 0
}

type SourceLocation struct {
	source string
	start  Position
	end    Position
}

//Node objects
type Node struct {
	jsT JsType
	loc SourceLocation
}

type Expression struct{ Node }

type Pattern struct{ Node }

type Statement struct{ Node }

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

type Program struct {
	Node
	body interface{} //body [ Directive | Statement ]

}
type Directive struct {
	expression Literal
	directive  string
}

type FunctionBody struct {
	BlockStatement
	body interface{}
	//body  [ Directive | Statement ]
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
