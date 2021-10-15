package parser

type SourceType string

const (
	ScriptType SourceType = "script"
	ModuleType SourceType = "module"
)

type ProgramES2015 struct {
	Program
	sourceType SourceType
	//body [ Statement | ModuleDeclaration ]
	//body interface{}
}

type FunctionES2015 struct {
	Function
	generator bool
}

type ForOfStatement struct {
	ForInStatement
}

const (
	LetType   KindType = "let"
	ConstType KindType = "const"
)

type Super struct {
	Node
}

type CallExpressionES2015 struct {
	CallExpression
	// callee: Expression | Super;
	callee interface{}
	//arguments: [ Expression | SpreadElement ]
	arguments interface{}
}

type MemberExpressionES2015 struct {
	MemberExpression
	//object: Expression | Super;
	object interface{}
}

type SpreadElement struct {
	Node
	argument Expression
}

type ArrayExpressionES2015 struct {
	ArrayExpression
	//elements: [ Expression | SpreadElement | null ];
	elements interface{}
}

type NewExpressionES2015 struct {
	NewExpression
	//arguments: [ Expression | SpreadElement ];
	arguments interface{}
}

type AssignmentExpressionES2015 struct {
	AssignmentExpression
	left Pattern
}

type PropertyES2015 struct {
	Property
	key       Expression
	method    bool
	shorthand bool
	computed  bool
}

type ArrowFunctionExpression struct {
	FunctionES2015
	Expression
	//body FunctionBody | Expression
	body       interface{}
	expression bool
}

type YieldExpression struct {
	Expression
	argument Expression
	delegate bool
}

type TemplateLiteral struct {
	Expression
	quasis      []TemplateElement
	expressions []Expression
}

type TaggedTemplateExpression struct {
	Expression
	tag   Expression
	quasi TemplateLiteral
}

type TemplateElement struct {
	Node
	tail  bool
	value struct {
		cooked string
		raw    string
	}
}

const InitKind KindType = "init"

type AssignmentProperty struct {
	PropertyES2015
	value  Pattern
	kind   KindType
	method bool
}

type ObjectPattern struct {
	Pattern
	properties []AssignmentProperty
}

type ArrayPattern struct {
	Pattern
	elements []Pattern
}

type RestElement struct {
	Pattern
	argument Pattern
}

type AssignmentPattern struct {
	Pattern
	left  Pattern
	right Expression
}

type Class struct {
	Node
	id         Identifier
	superClass Expression
	body       ClassBody
}

type ClassBody struct {
	Node
	body []MethodDefinition
}

const (
	ConstructorKind KindType = "constructor"
	MethodKind      KindType = "method"
	GetKind         KindType = "get"
	SetKind         KindType = "set"
)

type MethodDefinition struct {
	Node
	key      Expression
	value    FunctionExpression
	kind     KindType
	computed bool
	static   bool
}

type ClassDeclaration struct {
	Class
	Declaration
	id Identifier
}

type ClassExpression struct {
	Class
	Expression
}

type MetaProperty struct {
	Expression
	meta     Identifier
	property Identifier
}

type ModuleDeclaration struct {
	Node
}

type ModuleSpecifier struct {
	Node
	local Identifier
}

type ImportDeclaration struct {
	ModuleDeclaration
	//specifiers: [ ImportSpecifier | ImportDefaultSpecifier | ImportNamespaceSpecifier ];
	specifiers interface{}
	source     Literal
}

type ImportSpecifier struct {
	ModuleSpecifier
	imported Identifier
}

type ImportDefaultSpecifier struct {
	ModuleSpecifier
}

type ImportNamespaceSpecifier struct {
	ModuleSpecifier
}

type ExportNamedDeclaration struct {
	ModuleDeclaration
	declaration Declaration
	specifiers  []ExportSpecifier
	source      Literal
}

type ExportSpecifier struct {
	ModuleSpecifier
	exported Identifier
}

type AnonymousDefaultExportedFunctionDeclaration struct {
	Function
	//type: "FunctionDeclaration";
	//id: null;
}

type AnonymousDefaultExportedClassDeclaration struct {
	Class
	//type: "ClassDeclaration";
	//id: null;
}

type ExportDefaultDeclaration struct {
	ModuleDeclaration
	//declaration AnonymousDefaultExportedFunctionDeclaration | FunctionDeclaration | AnonymousDefaultExportedClassDeclaration | ClassDeclaration | Expression;
	declaration interface{}
}

type ExportAllDeclaration struct {
	ModuleDeclaration
	source Literal
}
