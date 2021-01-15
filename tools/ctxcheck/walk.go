package ctxcheck

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Walk(searchDir string, v *Visitor) error {
	return filepath.Walk(searchDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		if !strings.Contains(fi.Name(), ".go") {
			return nil
		}

		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		fs := token.NewFileSet()
		f, err := parser.ParseFile(fs, path, string(contents), parser.AllErrors)
		if err != nil {
			return err
		}

		v.SetPath(path)

		ast.Walk(v, f)

		if v.Err != nil {
			return v.Err
		}

		return nil
	})
}
