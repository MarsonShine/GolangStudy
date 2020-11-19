package eval

type Expr interface {
	Eval(env Env) float64
	// Check(vars map[Var]bool) error
}

type Var string
type literal float64

// unary和binary类型表示有一到两个运算对象的运算符表达式
type unary struct {
	op rune
	x  Expr
}
type binary struct {
	op   rune
	x, y Expr
}
type call struct {
	fn   string // 限制fn只能是pow，sin或者sqrt。
	args []Expr
}
