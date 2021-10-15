package parser

type ClassBodyES2022 struct {
	ClassBody
	//body: [ MethodDefinition | PropertyDefinition | StaticBlock ];
	body interface{}
}

type PropertyDefinition struct {
	Node
	//key: Expression | PrivateIdentifier;
	key      interface{}
	value    Expression
	computed bool
	static   bool
}

type MethodDefinitionES2022 struct {
	MethodDefinition
	//key: Expression | PrivateIdentifier;
	key interface{}
}

type PrivateIdentifier struct {
	Node
	name string
}

type MemberExpressionES2022 struct {
	MemberExpressionES2020
	//property: Expression | PrivateIdentifier;
	property interface{}
}

type StaticBlock struct {
	BlockStatement
}
type BinaryExpressionES2022 struct {
	Expression
	//left: Expression | PrivateIdentifier;
	left interface{}
}
