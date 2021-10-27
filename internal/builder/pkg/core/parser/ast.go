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
	Jsonify() NodeLike //to json
	Js() NodeLike      //to js
	Parse() NodeLike   //to ast
	isNodeLike() bool
}

type StatementLike interface {
	NodeLike
	isStatementLike() bool //just bind type , never call
}
type ExpressionLike interface {
	NodeLike
	isExpressionLike() bool //just bind type , never call
}

type PatternLike interface {
	NodeLike
	isPatternLike() bool //just bind type , never call
}

type FunctionLike interface {
	NodeLike
	isFunctionLike() bool //just bind type , never call
}

type DeclarationLike interface {
	StatementLike
	isDeclarationLike() bool //just bind type , never call
}

type LiteralLike interface {
	ExpressionLike
	isLiteralLike() bool //just bind type , never call
}

// BlockBodyLike 这个类型是一个特殊的类型，它是因为go 而特殊定制的
type BlockBodyLike interface {
	StatementLike
	isBlockBodyLike() bool //just bind type , never call
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
	Name []byte `json:"name"`
}

func (i Identifier) Jsonify() NodeLike {
	panic("implement me")
}

func (i Identifier) Js() NodeLike {
	panic("implement me")
}

func (i Identifier) Parse() NodeLike {
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
	//Value: Parse | boolean | null | number | RegExp;
	Value []byte `json:"value"`
	// 自带属性
	Regex struct {
		Pattern []byte `json:"pattern"`
		Flags   []byte `json:"flags"`
	} `json:"regex"`
}

func (r RegExpLiteral) Jsonify() NodeLike {
	panic("implement me")
}

func (r RegExpLiteral) Js() NodeLike {
	panic("implement me")
}

func (r RegExpLiteral) Parse() NodeLike {
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

func (p Programs) Parse() NodeLike {
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

func (f Function) Parse() NodeLike {
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

func (d Directive) Parse() NodeLike {
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

func (b BlockStatement) Parse() NodeLike {
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

func (e EmptyStatement) Parse() NodeLike {
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

func (d DebuggerStatement) Parse() NodeLike {
	panic("implement me")
}

func (d DebuggerStatement) isNodeLike() bool {
	panic("implement me")
}

func (d DebuggerStatement) isStatementLike() bool {
	panic("implement me")
}

//=============================================================================

type WithStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Object ExpressionLike `json:"object"`
	Body   StatementLike  `json:"body"`
}

func (w WithStatement) Jsonify() NodeLike {
	panic("implement me")
}

func (w WithStatement) Js() NodeLike {
	panic("implement me")
}

func (w WithStatement) Parse() NodeLike {
	panic("implement me")
}

func (w WithStatement) isNodeLike() bool {
	panic("implement me")
}

func (w WithStatement) isStatementLike() bool {
	panic("implement me")
}

//=============================================================================

type ReturnStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Argument ExpressionLike `json:"argument"`
}

func (r ReturnStatement) Jsonify() NodeLike {
	panic("implement me")
}

func (r ReturnStatement) Js() NodeLike {
	panic("implement me")
}

func (r ReturnStatement) Parse() NodeLike {
	panic("implement me")
}

func (r ReturnStatement) isNodeLike() bool {
	panic("implement me")
}

func (r ReturnStatement) isStatementLike() bool {
	panic("implement me")
}

//=============================================================================

type LabeledStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Label Identifier    `json:"label"`
	Body  StatementLike `json:"body"`
}

func (l LabeledStatement) Jsonify() NodeLike {
	panic("implement me")
}

func (l LabeledStatement) Js() NodeLike {
	panic("implement me")
}

func (l LabeledStatement) Parse() NodeLike {
	panic("implement me")
}

func (l LabeledStatement) isNodeLike() bool {
	panic("implement me")
}

func (l LabeledStatement) isStatementLike() bool {
	panic("implement me")
}

//=============================================================================

type BreakStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Label Identifier `json:"label"`
}

func (b BreakStatement) Jsonify() NodeLike {
	panic("implement me")
}

func (b BreakStatement) Js() NodeLike {
	panic("implement me")
}

func (b BreakStatement) Parse() NodeLike {
	panic("implement me")
}

func (b BreakStatement) isNodeLike() bool {
	panic("implement me")
}

func (b BreakStatement) isStatementLike() bool {
	panic("implement me")
}

//=============================================================================

type ContinueStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Label Identifier `json:"label"`
}

func (c ContinueStatement) Jsonify() NodeLike {
	panic("implement me")
}

func (c ContinueStatement) Js() NodeLike {
	panic("implement me")
}

func (c ContinueStatement) Parse() NodeLike {
	panic("implement me")
}

func (c ContinueStatement) isNodeLike() bool {
	panic("implement me")
}

func (c ContinueStatement) isStatementLike() bool {
	panic("implement me")
}

//=============================================================================

type IfStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Test       ExpressionLike `json:"test"`
	Consequent StatementLike  `json:"consequent"`
	Alternate  StatementLike  `json:"alternate"`
}

func (i IfStatement) Jsonify() NodeLike {
	panic("implement me")
}

func (i IfStatement) Js() NodeLike {
	panic("implement me")
}

func (i IfStatement) Parse() NodeLike {
	panic("implement me")
}

func (i IfStatement) isNodeLike() bool {
	panic("implement me")
}

func (i IfStatement) isStatementLike() bool {
	panic("implement me")
}

//=============================================================================
