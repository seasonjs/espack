package parser

type JsType string

// TODO 需要重新设计枚举，并且将lexer的枚举统一
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

// BlockBodyLike 这个类型是一个特殊的类型，它是因为go 而特殊定制的
type BlockBodyLike interface {
	StatementLike
	isBlockBodyLike() bool
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
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Name string `json:"name"`
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

type RegExpLiteral struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	//Literal
	//value: string | boolean | null | number | RegExp;
	value []byte
	// 自带属性
	Regex struct {
		Pattern string `json:"pattern"`
		Flags   string `json:"flags"`
	} `json:"regex"`
}

func (r RegExpLiteral) Jsonify() NodeLike {
	panic("implement me")
}

func (r RegExpLiteral) Js() NodeLike {
	panic("implement me")
}

func (r RegExpLiteral) String() NodeLike {
	panic("implement me")
}

func (r RegExpLiteral) isNodeLike() bool {
	panic("implement me")
}

func (r RegExpLiteral) isExpressionLike() bool {
	panic("implement me")
}

func (r RegExpLiteral) isLiteralLike() bool {
	panic("implement me")
}

//=============================================================================

type Programs struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Body []BlockBodyLike `json:"body"`
}

func (p Programs) Jsonify() NodeLike {
	panic("implement me")
}

func (p Programs) Js() NodeLike {
	panic("implement me")
}

func (p Programs) String() NodeLike {
	panic("implement me")
}

func (p Programs) isNodeLike() bool {
	panic("implement me")
}

//=============================================================================

type Function struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Id     Identifier    `json:"id"`
	Params []PatternLike `json:"params"`
	Body   BlockBodyLike `json:"body"`
}

func (f Function) Jsonify() NodeLike {
	panic("implement me")
}

func (f Function) Js() NodeLike {
	panic("implement me")
}

func (f Function) String() NodeLike {
	panic("implement me")
}

func (f Function) isNodeLike() bool {
	panic("implement me")
}

//=============================================================================

type Directive struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Expression LiteralLike `json:"expression"`
	Directive  []byte      `json:"directive"`
}

func (d Directive) Jsonify() NodeLike {
	panic("implement me")
}

func (d Directive) Js() NodeLike {
	panic("implement me")
}

func (d Directive) String() NodeLike {
	panic("implement me")
}

func (d Directive) isNodeLike() bool {
	panic("implement me")
}

//=============================================================================

type BlockStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Body []StatementLike `json:"body"`
}

func (b BlockStatement) Jsonify() NodeLike {
	panic("implement me")
}

func (b BlockStatement) Js() NodeLike {
	panic("implement me")
}

func (b BlockStatement) String() NodeLike {
	panic("implement me")
}

func (b BlockStatement) isNodeLike() bool {
	panic("implement me")
}

func (b BlockStatement) isStatementLike() bool {
	panic("implement me")
}

//=============================================================================

type EmptyStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
}

func (e EmptyStatement) Jsonify() NodeLike {
	panic("implement me")
}

func (e EmptyStatement) Js() NodeLike {
	panic("implement me")
}

func (e EmptyStatement) String() NodeLike {
	panic("implement me")
}

func (e EmptyStatement) isNodeLike() bool {
	panic("implement me")
}

func (e EmptyStatement) isStatementLike() bool {
	panic("implement me")
}

//=============================================================================

type DebuggerStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
}

func (d DebuggerStatement) Jsonify() NodeLike {
	panic("implement me")
}

func (d DebuggerStatement) Js() NodeLike {
	panic("implement me")
}

func (d DebuggerStatement) String() NodeLike {
	panic("implement me")
}

func (d DebuggerStatement) isNodeLike() bool {
	panic("implement me")
}

func (d DebuggerStatement) isStatementLike() bool {
	panic("implement me")
}

//=============================================================================
