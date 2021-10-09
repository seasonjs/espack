package core

type FunctionES2017 struct {
	FunctionES2015
	async bool
}

const AwaitExpressionType JsType = "AwaitExpression"

type AwaitExpression struct {
	Expression
	argument Expression
}
