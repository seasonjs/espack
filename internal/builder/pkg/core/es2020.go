package core

type LiteralES2020 struct {
	Literal
	value interface{}
}

type BigIntLiteral struct {
	Literal
	bigint string
}

const ChainExpressionType JsType = "ChainExpression"

type ChainExpression struct {
	Expression
	expression ChainElement
}

type ChainElement struct {
	Node
	optional bool
}

type CallExpressionES2020 struct {
	CallExpressionES2015
	ChainElement
}

type MemberExpressionES2020 struct {
	MemberExpressionES2015
	ChainElement
}

const ImportExpressionType JsType = "ImportExpression"

type ImportExpression struct {
	Expression
	source Expression
}

const LogicalOperatorNullishCoalescingOperator LogicalOperator = "??"

type ExportAllDeclarationES2020 struct {
	ExportAllDeclaration
	exported Identifier
}