package eval

import (
	"fmt"
	"math"
	"strings"
)

type Var string

type literal float64

type unary struct {
	op rune
	x  Expr
}

type binary struct {
	op   rune
	x, y Expr
}

type call struct {
	fn   string
	args []Expr
}

type Env map[Var]float64

type Expr interface {
	Eval(env Env) float64
	Check(vars map[Var]bool) error
}

func (v Var) Eval(env Env) float64 {
	return env[v]
}

func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

func (l literal) Check(vars map[Var]bool) error {
	return nil
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operarot: %q", u.op))
}

func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

var (
	add = func(x, y float64) float64 { return x + y }
	sub = func(x, y float64) float64 { return x - y }
	mul = func(x, y float64) float64 { return x * y }
	div = func(x, y float64) float64 { return x / y }
)

func (b binary) Eval(env Env) float64 {
	var f func(x, y float64) float64
	switch b.op {
	case '+':
		f = add
	case '-':
		f = sub
	case '*':
		f = mul
	case '/':
		f = div
	default:
		panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
	}
	return f(b.x.Eval(env), b.y.Eval(env))
}

func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/", u.op) {
		return fmt.Errorf("unexpected binary op %q", u.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

var (
	pow  = func(x, y float64) float64 { return math.Pow(x, y) }
	sin  = func(x float64) float64 { return math.Sin(x) }
	sqrt = func(x float64) float64 { return math.Sqrt(x) }
)

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return sin(c.args[0].Eval(env))
	case "sqrt":
		return sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

func (c call) Check(vars [Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d", c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}
