package js_parser

type JsType string

//需要注意的是这里的json转换是不符合我们需求的，所以应该直接使用string
// 下面是来自golang json包的转换文档
// String values encode as JSON strings coerced to valid UTF-8,
// replacing invalid bytes with the Unicode replacement rune.
// So that the JSON will be safe to embed inside HTML <script> tags,
// the string is encoded using HTMLEscape,
// which replaces "<", ">", "&", U+2028, and U+2029 are escaped
// to "\u003c","\u003e", "\u0026", "\u2028", and "\u2029".
// This replacement can be disabled when using an Encoder,
// by calling SetEscapeHTML(false).
//
// Array and slice values encode as JSON arrays, except that
// []byte encodes as a base64-encoded string, and a nil slice
// encodes as the null JSON value.

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
	isNodeLike()
}

type StatementLike interface {
	NodeLike
	isStatementLike() //just bind type , never call
}
type ExpressionLike interface {
	NodeLike
	isExpressionLike() //just bind type , never call
}

type PatternLike interface {
	NodeLike
	isPatternLike() //just bind type , never call
}

type FunctionLike interface {
	NodeLike
	isFunctionLike() //just bind type , never call
}

type DeclarationLike interface {
	StatementLike
	isDeclarationLike() //just bind type , never call
}

type LiteralLike interface {
	ExpressionLike
	isLiteralLike() //just bind type , never call
}

// BlockBodyLike 这个类型是一个特殊的类型，它是因为go 而特殊定制的
type BlockBodyLike interface {
	StatementLike
	isBlockBodyLike() //just bind type , never call
}

//=============================================================================

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type SourceLocation struct {
	Source string   `json:"source"`
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

func (i Identifier) Parse() NodeLike {
	panic("implement me")
}

func (i Identifier) isNodeLike() {
	panic("implement me")
}

func (i Identifier) isExpressionLike() {
	panic("implement me")
}

func (i Identifier) isPatternLike() {
	panic("implement me")
}

//=============================================================================

type RegExpLiteral struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	//Literal
	//Value: Parse | boolean | null | number | RegExp;
	Value string `json:"value"`
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

func (r RegExpLiteral) Parse() NodeLike {
	panic("implement me")
}

func (r RegExpLiteral) isNodeLike() {
	panic("implement me")
}

func (r RegExpLiteral) isExpressionLike() {
	panic("implement me")
}

func (r RegExpLiteral) isLiteralLike() {
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

func (p Programs) isNodeLike() {
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
	//es2015
	Generator bool `json:"generator"`
	//es2017
	Async bool `json:"async"`
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

func (f Function) isNodeLike() {
	panic("implement me")
}

//=============================================================================

type Directive struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Expression LiteralLike `json:"expression"`
	Directive  string      `json:"directive"`
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

func (d Directive) isNodeLike() {
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

func (b BlockStatement) isNodeLike() {
	panic("implement me")
}

func (b BlockStatement) isStatementLike() {
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

func (e EmptyStatement) isNodeLike() {
	panic("implement me")
}

func (e EmptyStatement) isStatementLike() {
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

func (d DebuggerStatement) isNodeLike() {
	panic("implement me")
}

func (d DebuggerStatement) isStatementLike() {
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

func (w WithStatement) isNodeLike() {
	panic("implement me")
}

func (w WithStatement) isStatementLike() {
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

func (r ReturnStatement) isNodeLike() {
	panic("implement me")
}

func (r ReturnStatement) isStatementLike() {
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

func (l LabeledStatement) isNodeLike() {
	panic("implement me")
}

func (l LabeledStatement) isStatementLike() {
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

func (b BreakStatement) isNodeLike() {
	panic("implement me")
}

func (b BreakStatement) isStatementLike() {
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

func (c ContinueStatement) isNodeLike() {
	panic("implement me")
}

func (c ContinueStatement) isStatementLike() {
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

func (i IfStatement) isNodeLike() {
	panic("implement me")
}

func (i IfStatement) isStatementLike() {
	panic("implement me")
}

//=============================================================================

type SwitchCase struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Test       ExpressionLike  `json:"test"`
	Consequent []StatementLike `json:"consequent"`
}

func (s SwitchCase) Jsonify() NodeLike {
	panic("implement me")
}

func (s SwitchCase) Js() NodeLike {
	panic("implement me")
}

func (s SwitchCase) Parse() NodeLike {
	panic("implement me")
}

func (s SwitchCase) isNodeLike() {
	panic("implement me")
}

//=============================================================================

type SwitchStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Discriminant ExpressionLike `json:"discriminant"`
	Cases        []SwitchCase   `json:"cases"`
}

func (s SwitchStatement) Jsonify() NodeLike {
	panic("implement me")
}

func (s SwitchStatement) Js() NodeLike {
	panic("implement me")
}

func (s SwitchStatement) Parse() NodeLike {
	panic("implement me")
}

func (s SwitchStatement) isNodeLike() {
	panic("implement me")
}

func (s SwitchStatement) isStatementLike() {
	panic("implement me")
}

//=============================================================================

type ThrowStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Argument ExpressionLike `json:"argument"`
}

func (t ThrowStatement) Jsonify() NodeLike {
	panic("implement me")
}

func (t ThrowStatement) Js() NodeLike {
	panic("implement me")
}

func (t ThrowStatement) Parse() NodeLike {
	panic("implement me")
}

func (t ThrowStatement) isNodeLike() {
	panic("implement me")
}

func (t ThrowStatement) isStatementLike() {
	panic("implement me")
}

//=============================================================================

type TryStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Block     BlockStatement `json:"block"`
	Handler   CatchClause    `json:"handler"`
	Finalizer BlockStatement `json:"finalizer"`
}

func (t TryStatement) Jsonify() NodeLike {
	panic("implement me")
}

func (t TryStatement) Js() NodeLike {
	panic("implement me")
}

func (t TryStatement) Parse() NodeLike {
	panic("implement me")
}

func (t TryStatement) isNodeLike() {
	panic("implement me")
}

func (t TryStatement) isStatementLike() {
	panic("implement me")
}

//=============================================================================

type CatchClause struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Param PatternLike    `json:"param"` // es2019 The param is null if the catch binding is omitted. E.g., try { foo() } catch { bar() }
	Body  BlockStatement `json:"body"`
}

func (c CatchClause) Jsonify() NodeLike {
	panic("implement me")
}

func (c CatchClause) Js() NodeLike {
	panic("implement me")
}

func (c CatchClause) Parse() NodeLike {
	panic("implement me")
}

func (c CatchClause) isNodeLike() {
	panic("implement me")
}

//=============================================================================

type WhileStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Test ExpressionLike `json:"test"`
	Body StatementLike  `json:"body"`
}

//=============================================================================

type DoWhileStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Body StatementLike  `json:"body"`
	Test ExpressionLike `json:"test"`
}

//=============================================================================

type ForInitBodyLike interface {
	StatementLike
	DeclarationLike
}

type ForStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Init   ForInitBodyLike `json:"init"`
	Test   ExpressionLike  `json:"test"`
	Update ExpressionLike  `json:"update"`
	Body   StatementLike   `json:"body"`
}

//=============================================================================

type ForLeftBodyLike interface {
	NodeLike
}

type ForInStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Left  ForLeftBodyLike `json:"left"`
	Right ExpressionLike  `json:"right"`
	Body  StatementLike   `json:"body"`
}

//=============================================================================

type VariableDeclaration struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Declarations []VariableDeclarator `json:"declarations"`
	Kind         string               `json:"kind"` //es5  "var"// es2015 "var" | "let" | "const" 这里是可以使用string的
}

//=============================================================================

type VariableDeclarator struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Id   PatternLike    `json:"id"`
	Init ExpressionLike `json:"init"`
}

//=============================================================================

type ThisExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
}

//=============================================================================

type ArrayExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Elements []ExpressionLike `json:"elements"` //es2015 Spread element
}

//=============================================================================

type ObjectExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Properties []Property `json:"properties"` //TODO  properties: [ Property | SpreadElement ];
}

//=============================================================================

type PropertyKeyBodyLike interface {
	NodeLike
}

type Property struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Key   PropertyKeyBodyLike `json:"key"`
	Value ExpressionLike      `json:"value"`
	Kind  string              `json:"kind"`
	//es2015
	Method    bool `json:"method"`
	Shorthand bool `json:"shorthand"`
	Computed  bool `json:"computed"`
}

//=============================================================================

type FunctionExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Id     Identifier    `json:"id"`
	Params []PatternLike `json:"params"`
	Body   BlockBodyLike `json:"body"`
	//es2015
	Generator bool `json:"generator"`
	//es2017
	Async bool `json:"async"`
}

//=============================================================================

type BinaryOperator string

type BinaryExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Operator BinaryOperator `json:"operator"`
	Left     ExpressionLike `json:"left"` //es2020 ExpressionLike | PrivateIdentifier;
	Right    ExpressionLike `json:"right"`
}

//=============================================================================

type AssignmentOperator string

type AssignmentLeftLike interface {
	NodeLike
}

type AssignmentExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Operator AssignmentOperator `json:"operator"`
	Left     AssignmentLeftLike `json:"left"`
	Right    ExpressionLike     `json:"right"`
}

//=============================================================================

type LogicalOperator string

type LogicalExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Operator LogicalOperator `json:"operator"`
	Left     ExpressionLike  `json:"left"`
	Right    ExpressionLike  `json:"right"`
}

//=============================================================================

type MemberExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Object   ExpressionLike `json:"object"`   //es2015 super also is
	Property ExpressionLike `json:"property"` //es2022 Expression | PrivateIdentifier;
	Computed bool           `json:"computed"`
	Optional bool           `json:"optional"` //es2020 ChainElement
}

//=============================================================================

type ConditionalExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Test       ExpressionLike `json:"test"`
	Alternate  ExpressionLike `json:"alternate"`
	Consequent ExpressionLike `json:"consequent"`
}

//=============================================================================

type CallExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Callee    ExpressionLike   `json:"callee"`    //es2015 Super
	Arguments []ExpressionLike `json:"arguments"` //es2015 SpreadElement
	Optional  bool             `json:"optional"`  //es2021 ChainElement
}

//=============================================================================

type NewExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Callee    ExpressionLike   `json:"callee"`
	Arguments []ExpressionLike `json:"arguments"` //es2015 SpreadElement
}

//=============================================================================

type SequenceExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Expressions []ExpressionLike `json:"expressions"`
}

//es2015
//=============================================================================

type Program struct {
	sourceType string
	body       []BlockBodyLike
}

//=============================================================================

type ForOfStatement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	Left  ForLeftBodyLike `json:"left"`
	Right ExpressionLike  `json:"right"`
	Body  StatementLike   `json:"body"`
	//es2018
	Await bool `json:"await"`
}

//=============================================================================

type Super struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
}

func (s Super) Jsonify() NodeLike {
	panic("implement me")
}

func (s Super) Js() NodeLike {
	panic("implement me")
}

func (s Super) Parse() NodeLike {
	panic("implement me")
}

func (s Super) isNodeLike() {
	panic("implement me")
}

//=============================================================================

type SpreadElement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
	// 自带属性
	argument ExpressionLike
}

//=============================================================================

type ArrowFunctionExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Id     Identifier    `json:"id"`
	Params []PatternLike `json:"params"`
	Body   BlockBodyLike `json:"body"`
	//es2015
	Generator  bool `json:"generator"` // 必须是false
	Expression bool `json:"expression"`
	//es2017
	Async bool `json:"async"`
}

//=============================================================================

type YieldExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Argument ExpressionLike `json:"argument"`
	Delegate bool           `json:"delegate"`
}

//=============================================================================

type TemplateElement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Tail  bool `json:"tail"`
	Value struct {
		Cooked string `json:"cooked"` //in es2018 If the template literal is tagged and the text has an invalid escape, cooked will be null, e.g., tag`\unicode and \u{55}`
		Raw    string `json:"raw"`
	} `json:"value"`
}

//=============================================================================

type TemplateLiteral struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	quasis      []TemplateElement `json:"quasis"`
	expressions []ExpressionLike  `json:"expressions"`
}

//=============================================================================

type AssignmentProperty struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Value     PatternLike         `json:"value"`
	Kind      string              `json:"kind"`
	Method    bool                `json:"method"` //must be false
	Key       PropertyKeyBodyLike `json:"key"`
	Shorthand bool                `json:"shorthand"`
	Computed  bool                `json:"computed"`
}

//=============================================================================

type ObjectPattern struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Properties []AssignmentProperty `json:"properties"` // TODO in es2018 properties: [ AssignmentProperty | RestElement ];
}

//=============================================================================

type ArrayPattern struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Elements []PatternLike `json:"elements"`
}

//=============================================================================

type RestElement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Argument PatternLike `json:"argument"`
}

//=============================================================================

type AssignmentPattern struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Left  PatternLike    `json:"left"`
	Right ExpressionLike `json:"right"`
}

//=============================================================================

type Class struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	superClass ExpressionLike
	body       ClassBody
}

//=============================================================================

type ClassBodyLike interface {
	isClassBodyLike() bool
}

//=============================================================================

type ClassBody struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Body []ClassBodyLike `json:"body"`
}

//=============================================================================

type MethodDefinition struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Key      ExpressionLike     `json:"key"` // es2021  key: Expression | PrivateIdentifier;
	Value    FunctionExpression `json:"value"`
	Kind     string             `json:"kind"`
	Computed bool               `json:"computed"`
	Static   bool               `json:"static"`
}

//=============================================================================

type ClassDeclaration struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Id Identifier `json:"id"`
}

//=============================================================================

type ClassExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	superClass ExpressionLike
	body       ClassBody
}

//=============================================================================

type MetaProperty struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Meta     Identifier `json:"meta"`
	Property Identifier `json:"property"`
}

//=============================================================================

type ModuleDeclaration struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`
}

func (m ModuleDeclaration) Jsonify() NodeLike {
	panic("implement me")
}

func (m ModuleDeclaration) Js() NodeLike {
	panic("implement me")
}

func (m ModuleDeclaration) Parse() NodeLike {
	panic("implement me")
}

func (m ModuleDeclaration) isNodeLike() {
	panic("implement me")
}

//=============================================================================

type ModuleSpecifier struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	Local Identifier `json:"local"`
}

func (m ModuleSpecifier) Jsonify() NodeLike {
	panic("implement me")
}

func (m ModuleSpecifier) Js() NodeLike {
	panic("implement me")
}

func (m ModuleSpecifier) Parse() NodeLike {
	panic("implement me")
}

func (m ModuleSpecifier) isNodeLike() {
	panic("implement me")
}

//=============================================================================

type ImportSpecifierLike interface {
	NodeLike
	isImportSpecifierLike() bool
}

//=============================================================================

type ImportDeclaration struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Specifiers []ImportSpecifierLike `json:"specifiers"`
	Source     LiteralLike           `json:"source"`
}

//=============================================================================

type ImportSpecifier struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	Local    Identifier `json:"local"`
	Imported Identifier `json:"imported"`
}

func (i ImportSpecifier) Jsonify() NodeLike {
	panic("implement me")
}

func (i ImportSpecifier) Js() NodeLike {
	panic("implement me")
}

func (i ImportSpecifier) Parse() NodeLike {
	panic("implement me")
}

func (i ImportSpecifier) isNodeLike() {
	panic("implement me")
}

func (i ImportSpecifier) isImportSpecifierLike() bool {
	panic("implement me")
}

//=============================================================================

type ImportDefaultSpecifier struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	Local Identifier `json:"local"`
}

func (i ImportDefaultSpecifier) Jsonify() NodeLike {
	panic("implement me")
}

func (i ImportDefaultSpecifier) Js() NodeLike {
	panic("implement me")
}

func (i ImportDefaultSpecifier) Parse() NodeLike {
	panic("implement me")
}

func (i ImportDefaultSpecifier) isNodeLike() {
	panic("implement me")
}

func (i ImportDefaultSpecifier) isImportSpecifierLike() bool {
	panic("implement me")
}

//=============================================================================

type ImportNamespaceSpecifier struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	Local Identifier `json:"local"`
}

func (i ImportNamespaceSpecifier) Jsonify() NodeLike {
	panic("implement me")
}

func (i ImportNamespaceSpecifier) Js() NodeLike {
	panic("implement me")
}

func (i ImportNamespaceSpecifier) Parse() NodeLike {
	panic("implement me")
}

func (i ImportNamespaceSpecifier) isNodeLike() {
	panic("implement me")
}

func (i ImportNamespaceSpecifier) isImportSpecifierLike() bool {
	panic("implement me")
}

//=============================================================================

type ExportNamedDeclaration struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	Local       Identifier        `json:"local"`
	Declaration DeclarationLike   `json:"declaration"`
	Specifiers  []ExportSpecifier `json:"specifiers"`
	Source      LiteralLike       `json:"source"`
}

//=============================================================================

type ExportSpecifier struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	Local    Identifier `json:"local"`
	Exported Identifier `json:"exported"`
}

//=============================================================================

type AnonymousDefaultExportedFunctionDeclaration struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Id     Identifier    `json:"id"` //must be null
	Params []PatternLike `json:"params"`
	Body   BlockBodyLike `json:"body"`

	//es2015
	Generator bool `json:"generator"`

	//es2017
	Async bool `json:"async"`
}

//=============================================================================

type AnonymousDefaultExportedClassDeclaration struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Id         Identifier `json:"id"` //must be null
	superClass ExpressionLike
	body       ClassBody
}

//=============================================================================

type ExportDefaultDeclarationLike interface {
	isExportDefaultDeclarationLike() bool
}

//=============================================================================

type ExportDefaultDeclaration struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	Local       Identifier                   `json:"local"`
	Declaration ExportDefaultDeclarationLike `json:"declaration"`
}

//=============================================================================

type ExportAllDeclaration struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	Local    Identifier  `json:"local"`
	Source   LiteralLike `json:"source"`
	Exported Identifier  `json:"exported"` //es2020
}

//es2020
//=============================================================================

type BigIntLiteral struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	Bigint string `json:"bigint"`
}

//=============================================================================

type ChainExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	expression ChainElement
}

//=============================================================================

type ChainElement struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	Optional bool `json:"optional"`
}

//=============================================================================

type ImportExpression struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	source ExpressionLike
}

//es2022
//=============================================================================

type PropertyDefinition struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	Key      ExpressionLike `json:"key"`
	Value    ExpressionLike `json:"value"`
	Computed bool           `json:"computed"`
	Static   bool           `json:"static"`
}

//=============================================================================

type PrivateIdentifier struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	//自带属性
	Name string `json:"name"`
}

//=============================================================================

type StaticBlock struct {
	// 从Node 节点继承
	Loc SourceLocation `json:"loc"`
	JsT JsType         `json:"type"`

	// 自带属性
	Body []StatementLike `json:"body"`
}
