package parser

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

//=============================================================================
//最顶层的ast结构为了实现 go duck type

type NodeLike interface {
	Jsonify() NodeLike
	Js() NodeLike
	String() NodeLike
	isNodeLike() bool
}

type StatementLike interface {
	NodeLike
	isStatementLike() bool
}
type ExpressionLike interface {
	NodeLike
	isExpressionLike() bool
}

type PatternLike interface {
	NodeLike
	isPatternLike() bool
}

type FunctionLike interface {
	NodeLike
	isFunctionLike() bool
}

type DeclarationLike interface {
	StatementLike
	isDeclarationLike() bool
}

type LiteralLike interface {
	ExpressionLike
	isLiteralLike() bool
}

//=============================================================================

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type SourceLocation struct {
	Source []byte   `json:"source"`
	Start  Position `json:"start"`
	End    Position `json:"end"`
}

//=============================================================================

type Identifier struct {
	Loc  SourceLocation `json:"loc"`
	JsT  JsType         `json:"type"`
	Name string         `json:"name"`
}

func (i Identifier) Jsonify() NodeLike {
	panic("implement me")
}

func (i Identifier) Js() NodeLike {
	panic("implement me")
}

func (i Identifier) String() NodeLike {
	panic("implement me")
}

func (i Identifier) isNodeLike() bool {
	panic("implement me")
}

func (i Identifier) isExpressionLike() bool {
	panic("implement me")
}

func (i Identifier) isPatternLike() bool {
	panic("implement me")
}

//=============================================================================
