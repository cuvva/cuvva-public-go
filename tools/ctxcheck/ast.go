package ctxcheck

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
)

type FunctionPredicate func(i *ast.Ident) bool
type ArgPredicate func(index int, ident *ast.Ident) *ast.Ident

type Visitor struct {
	path  string
	found bool
	Err   error

	functionPredicate FunctionPredicate
	argPredicate      ArgPredicate

	Matches []*CheckerMatch
}

func NewVisitor(fnPredicate FunctionPredicate, argPredicate ArgPredicate) *Visitor {
	return &Visitor{
		found:             false,
		functionPredicate: fnPredicate,
		argPredicate:      argPredicate,
	}
}

func (v *Visitor) Visit(n ast.Node) (next ast.Visitor) {
	if n == nil {
		return nil
	}

	if v.found {
		v.found = false

		var match *CheckerMatch
		next, match, v.Err = v.findArgInFunction(n)

		if match != nil {
			v.Matches = append(v.Matches, match)
		}

		return
	}

	v.found, next = v.findIdentInNode(n)

	return
}

func (v *Visitor) SetPath(path string) {
	v.path = path
}

func (v *Visitor) findArgInFunction(n ast.Node) (ast.Visitor, *CheckerMatch, error) {
	funcLit, ok := n.(*ast.FuncLit)
	if !ok || len(funcLit.Body.List) == 0 {
		return v, nil, nil
	}

	sc := statementChecker{
		predicate: v.argPredicate,
		path:      v.path,
	}

	err := sc.handleStmt(funcLit.Body.List[0])
	if err != nil {
		return nil, nil, err
	}

	return nil, sc.match, nil
}

type statementChecker struct {
	predicate func(index int, ident *ast.Ident) *ast.Ident
	path      string
	match     *CheckerMatch
}

type CheckerMatch struct {
	Path string
	Pos  int
}

func NewCheckerMatch(path string, pos token.Pos) *CheckerMatch {
	return &CheckerMatch{
		Path: path,
		Pos:  int(pos),
	}
}

func (sc *statementChecker) handleStmt(stmt ast.Stmt) (err error) {
	if stmt == nil {
		return nil
	}

	switch stmt.(type) {
	case *ast.AssignStmt:
		err = sc.handleAssignStmt(stmt)
	case *ast.ExprStmt:
		err = sc.handleExprStmt(stmt)
	case *ast.ReturnStmt:
		err = sc.handleReturnStmt(stmt)
	case *ast.IfStmt:
		err = sc.handleIfStmt(stmt)
	case *ast.RangeStmt:
	default:
		return fmt.Errorf("not yet supported type %T", stmt)
	}

	return
}

func (sc *statementChecker) handleAssignStmt(raw ast.Stmt) error {
	stmt, ok := raw.(*ast.AssignStmt)
	if !ok {
		return errors.New("incorrect stmt provided")
	}

	for _, expr := range stmt.Rhs {
		call, ok := expr.(*ast.CallExpr)
		if !ok {
			continue
		}

		foundArg := findInArgs(call.Args, sc.predicate)
		if foundArg != nil {
			sc.match = NewCheckerMatch(sc.path, foundArg.Pos())
		}
	}

	return nil
}

func (sc *statementChecker) handleExprStmt(raw ast.Stmt) error {
	stmt, ok := raw.(*ast.ExprStmt)
	if !ok {
		return errors.New("incorrect stmt provided")
	}

	call, ok := stmt.X.(*ast.CallExpr)
	if !ok {
		return errors.New("incorrect stmt provided")
	}

	foundArg := findInArgs(call.Args, sc.predicate)
	if foundArg != nil {
		sc.match = NewCheckerMatch(sc.path, foundArg.Pos())
	}

	return nil
}

func (sc *statementChecker) handleReturnStmt(raw ast.Stmt) error {
	stmt, ok := raw.(*ast.ReturnStmt)
	if !ok {
		return errors.New("incorrect stmt provided")
	}

	for _, res := range stmt.Results {
		call, ok := res.(*ast.CallExpr)
		if !ok {
			return nil
		}

		foundArg := findInArgs(call.Args, sc.predicate)
		if foundArg != nil {
			sc.match = NewCheckerMatch(sc.path, foundArg.Pos())
		}
	}

	return nil
}

func (sc *statementChecker) handleIfStmt(raw ast.Stmt) error {
	stmt, ok := raw.(*ast.IfStmt)
	if !ok {
		return errors.New("incorrect stmt provided")
	}

	return sc.handleStmt(stmt.Init)
}

func findInArgs(args []ast.Expr, pred func(i int, ident *ast.Ident) *ast.Ident) *ast.Ident {
	for i, arg := range args {
		ident, ok := arg.(*ast.Ident)
		if !ok {
			continue
		}

		found := pred(i, ident)
		if found != nil {
			return found
		}
	}

	return nil
}

func (v *Visitor) findIdentInNode(n ast.Node) (bool, ast.Visitor) {
	ident, ok := n.(*ast.Ident)
	if !ok {
		return false, v
	}

	if v.functionPredicate(ident) {
		return true, v
	}

	return false, v
}
