package aggregation

const (
	NOW Expression = varExpression("NOW")
)

//Expression can contain a Field, Literal, SysVar,
//ExpressionMap or Operator
type Expression interface {
	Apply() interface{}
}

type SysVar string

func (s SysVar) Apply() interface{} {
	return "$$" + f
}

//Field accesses a field in the input documents.
func Field(name string) Expression {
	return fieldExpression(name)
}

//Var can be a user-defined or system variable.
func Var(name string) Expression {
	return varExpression(name)
}

type fieldExpression string
type varExpression string

func (f fieldExpression) Apply() interface{} {
	return "$" + f
}

func (v varExpression) Apply() interface{} {
	return "$$" + v
}
