package aggregation

//Expression can contain a Field, Literal, SysVar,
//ExpressionMap or Operator
type Expression interface {
	Apply() interface{}
}
