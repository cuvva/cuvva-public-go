package ctxcheck_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/cuvva/cuvva-public-go/tools/ctxcheck"
)

const exampleCode = `
package main

func main() {
	ctx := context.Background()

	var g errgroup.Group

	g.Go(func (gctx context.Context) error {
		requestSomeService(ctx)
	})

	if strings.Contains("a", "ab") {
		fmt.Println("hello")
	}

	requestSomeService(ctx)

	Go(ctx)
}

func Go(ctx context.Context) {
	return
}

func requestSomeService(_ context.Context) {
	return
}
`

func TestCtxFind(t *testing.T) {
	var functionPredicate ctxcheck.FunctionPredicate = func(i *ast.Ident) bool {
		return i.Name == "Go"
	}

	var argumentPredicate ctxcheck.ArgPredicate = func(index int, ident *ast.Ident) *ast.Ident {
		if index == 0 && ident.Name == "ctx" {
			return ident
		}

		return nil
	}

	v := ctxcheck.NewVisitor(functionPredicate, argumentPredicate)

	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "local", exampleCode, parser.AllErrors)
	if err != nil {
		t.Error(err)
	}

	ast.Walk(v, f)

	if v.Err != nil {
		t.Error(v.Err)
	}

	if len(v.Matches) != 1 {
		t.Errorf("expected 1 match, got %d", len(v.Matches))
	}

	if v.Matches[0].Pos != 146 {
		t.Errorf("expected to find a bad context as pos 146, found at %d", v.Matches[0].Pos)
	}
}
