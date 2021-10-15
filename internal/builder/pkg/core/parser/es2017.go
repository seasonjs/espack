package parser

type FunctionES2017 struct {
	FunctionES2015
	async bool
}

type AwaitExpression struct {
	Expression
	argument Expression
}
