package cleaner

import (
	"go/ast"
	"go/token"
	"log/slog"

	"github.com/Fabianexe/gocoverageplus/pkg/entity"
)

// cleanErrorIf removes all error if statements from the package data
// An error if statement is an if statement that checks if a variable x is not nil, x is named err or is of type error and has only a return statement in the body.
func cleanErrorIf(project *entity.Project) *entity.Project {
	slog.Info("Clean error if statements")
	for _, p := range project.Packages {
		for _, f := range p.Files {
			var countErrorIf int
			for _, method := range f.Methods {
				cleanErrorIfVisitor := &cleanErrorIfVisitor{
					fset: p.Fset,
				}

				ast.Walk(cleanErrorIfVisitor, method.Body)

				countErrorIf += len(cleanErrorIfVisitor.errorIfs)

				for _, errIf := range cleanErrorIfVisitor.errorIfs {
					method.Tree.AddBlock(
						&entity.Block{
							StartPosition: errIf.start,
							EndPosition:   errIf.end,
							DefPosition:   errIf.start,
							Type:          entity.TypeBlock,
							Ignore:        true,
						},
					)
				}
			}
			if countErrorIf > 0 {
				slog.Debug("Cleaned error if statements", "File", f.FilePath, "ErrorIfs", countErrorIf)
			}
		}
	}

	return project
}

type cleanErrorIfVisitor struct {
	errorIfs []errorIF
	fset     *token.FileSet
}

type errorIF struct {
	start token.Position
	end   token.Position
}

func (c *cleanErrorIfVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if !isErrorIf(node) {
		return c
	}

	c.errorIfs = append(c.errorIfs, errorIF{
		start: c.fset.Position(node.Pos()),
		end:   c.fset.Position(node.End()),
	})

	return c
}

func isErrorIf(node ast.Node) bool {
	if node == nil { //  nothing here
		return false
	}

	v, ok := node.(*ast.IfStmt)
	if !ok { // no if statement
		return false
	}

	if v.Cond == nil || // np condition
		v.Else != nil || // has else
		len(v.Body.List) != 1 { // more than one statement in Body
		return false
	}

	cond, ok := v.Cond.(*ast.BinaryExpr)
	if !ok { // no binary expression
		return false
	}

	if cond.Op != token.NEQ { // not !=
		return false
	}

	if compare, ok := cond.Y.(*ast.Ident); !ok || compare.Name != "nil" { // not compared against nil
		return false
	}

	if _, ok := v.Body.List[0].(*ast.ReturnStmt); !ok { // body is not a return statement
		return false
	}

	return isErrorVar(cond)
}

func isErrorVar(cond *ast.BinaryExpr) bool {
	object, ok := cond.X.(*ast.Ident)
	if !ok || object.Obj == nil || object.Obj.Kind != ast.Var { // not a variable
		return false
	}

	if object.Name != "err" { // not named err
		// try to determine type of object
		typ, ok := object.Obj.Decl.(*ast.ValueSpec)
		if !ok { // not a value spec
			return false
		}

		if n, ok := typ.Type.(*ast.Ident); !ok || n.Name != "error" { // not an error
			return false
		}
	}

	return true
}
