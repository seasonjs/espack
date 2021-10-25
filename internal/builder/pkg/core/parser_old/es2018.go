package parser_old

type ForOfStatementES2017 struct {
	ForOfStatement
	await bool
}
type ObjectExpressionES2017 struct {
	ObjectExpression
	//properties [ PropertyES2015 | SpreadElement ];
	Property interface{}
}
type TemplateElementES2017 struct {
	TemplateElement
	value struct {
		cooked string
		raw    string
	}
}
type ObjectPatternES2017 struct {
	ObjectPattern
	//properties: [ AssignmentProperty | RestElement ];
	properties interface{}
}
