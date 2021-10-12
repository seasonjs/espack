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

const SuperType JsType = "super"

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

const SpreadElementType JsType = "SpreadElement"

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

const ArrowFunctionExpressionType JsType = "ArrowFunctionExpression"

type ArrowFunctionExpression struct {
	FunctionES2015
	Expression
	//body FunctionBody | Expression
	body       interface{}
	expression bool
}

const YieldExpressionType JsType = "YieldExpression"

type YieldExpression struct {
	Expression
	argument Expression
	delegate bool
}

const ExpressionType JsType = "Expression"

type TemplateLiteral struct {
	Expression
	quasis      []TemplateElement
	expressions []Expression
}

const TaggedTemplateExpressionType JsType = "TaggedTemplateExpression"

type TaggedTemplateExpression struct {
	Expression
	tag   Expression
	quasi TemplateLiteral
}

const TemplateElementType JsType = "TemplateElement"

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

const ObjectPatternType = "ObjectPattern"

type ObjectPattern struct {
	Pattern
	properties []AssignmentProperty
}

const ArrayPatternType JsType = "ArrayPattern"

type ArrayPattern struct {
	Pattern
	elements []Pattern
}

const RestElementType JsType = "RestElement"

type RestElement struct {
	Pattern
	argument Pattern
}

const AssignmentPatternType JsType = "AssignmentPattern"

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

const ClassBodyType JsType = "ClassBody"

type ClassBody struct {
	Node
	body []MethodDefinition
}

const MethodDefinitionType JsType = "MethodDefinition"

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

const ClassDeclarationType JsType = "ClassDeclaration"

type ClassDeclaration struct {
	Class
	Declaration
	id Identifier
}

const ClassExpressionType JsType = "ClassExpression"

type ClassExpression struct {
	Class
	Expression
}

const MetaPropertyType JsType = "MetaProperty"

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

const ImportDeclarationType JsType = "ImportDeclaration"

type ImportDeclaration struct {
	ModuleDeclaration
	//specifiers: [ ImportSpecifier | ImportDefaultSpecifier | ImportNamespaceSpecifier ];
	specifiers interface{}
	source     Literal
}

const ImportSpecifierType = "ImportSpecifier"

type ImportSpecifier struct {
	ModuleSpecifier
	imported Identifier
}

const ImportDefaultSpecifierType JsType = "ImportDefaultSpecifier"

type ImportDefaultSpecifier struct {
	ModuleSpecifier
}

const ImportNamespaceSpecifierType JsType = "ImportNamespaceSpecifier"

type ImportNamespaceSpecifier struct {
	ModuleSpecifier
}

const ExportNamedDeclarationType JsType = "ExportNamedDeclaration"

type ExportNamedDeclaration struct {
	ModuleDeclaration
	declaration Declaration
	specifiers  []ExportSpecifier
	source      Literal
}

const ExportSpecifierType JsType = "ExportSpecifier"

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

const ExportDefaultDeclarationType JsType = "ExportDefaultDeclaration"

type ExportDefaultDeclaration struct {
	ModuleDeclaration
	//declaration AnonymousDefaultExportedFunctionDeclaration | FunctionDeclaration | AnonymousDefaultExportedClassDeclaration | ClassDeclaration | Expression;
	declaration interface{}
}

const ExportAllDeclarationType JsType = "ExportAllDeclaration"

type ExportAllDeclaration struct {
	ModuleDeclaration
	source Literal
}
