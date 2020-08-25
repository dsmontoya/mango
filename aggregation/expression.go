package aggregation

import "github.com/dsmontoya/mango/bson"

const (
	NOW = varExpression("NOW")
)

//Expression can contain a Field, Literal, SysVar,
//ExpressionMap or Operator
type Expression interface {
	Apply() interface{}
}

//ExpressionObject represents the
//{ <field>: <expression1>, ... } form.
type ExpressionObject map[string]Expression

type fieldExpression string
type varExpression string

type opExpression struct {
	apply func() interface{}
}

//Apply returns an expression object.
func (e ExpressionObject) Apply() interface{} {
	return e
}

//Apply returns an expression object.
func (o *opExpression) Apply() interface{} {
	return o.apply()
}

//Field accesses a field in the input documents.
func Field(name string) Expression {
	return fieldExpression(name)
}

//SetUnion takes two or more arrays and returns an array
//containing the elements that appear in any input array.
//The arguments can be any Expression as long as they
//each resolve to an array.
//
//See https://docs.mongodb.com/manual/reference/operator/aggregation/setUnion/#exp._S_setUnion.
func SetUnion(expressions []Expression) Expression {
	return &opExpression{
		apply: func() interface{} {
			exps := make([]interface{}, len(expressions))
			for i, expression := range expressions {
				exps[i] = expression.Apply()
			}
			return bson.M{"$setUnion": exps}
		},
	}
}

//Var can be a user-defined or system variable.
func Var(name string) Expression {
	return varExpression(name)
}

func (f fieldExpression) Apply() interface{} {
	return string("$" + f)
}

func (v varExpression) Apply() interface{} {
	return string("$$" + v)
}
