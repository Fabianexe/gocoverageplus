package cleaner

import (
	"go/ast"
	"go/token"
	"log/slog"
	"slices"

	"github.com/Fabianexe/gocoverageplus/pkg/entity"
)

// cleanErrorIf removes all error if statements from the package data
// An error if statement is an if statement that checks if a variable x is not nil, x is named err or is of type error and has only a return statement in the body.
func cleanCustomIf(project *entity.Project, names []string) *entity.Project {
	slog.Info("Clean error if statements")
	for _, p := range project.Packages {
		for _, f := range p.Files {
			var countCustomIf int
			for _, method := range f.Methods {
				vis := &cleanCustomIfVisitor{
					fset:  p.Fset,
					names: names,
				}

				ast.Walk(vis, method.Body)

				countCustomIf += len(vis.customIF)

				for _, errIf := range vis.customIF {
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
			if countCustomIf > 0 {
				slog.Debug("Cleaned custom if statements", "File", f.FilePath, "ErrorIfs", countCustomIf)
			}
		}
	}

	return project
}

type cleanCustomIfVisitor struct {
	customIF []customIF
	fset     *token.FileSet
	names    []string
}

type customIF struct {
	start token.Position
	end   token.Position
}

func (c *cleanCustomIfVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if !checkCustomIf(node, c.names) {
		return c
	}

	c.customIF = append(c.customIF, customIF{
		start: c.fset.Position(node.(*ast.IfStmt).Cond.Pos()),
		end:   c.fset.Position(node.End()),
	})

	return c
}

func checkCustomIf(node ast.Node, names []string) bool {
	if node == nil { //  nothing here
		return false
	}

	v, ok := node.(*ast.IfStmt)
	if !ok { // no if statement
		return false
	}

	if v.Cond == nil || // no condition
		v.Else != nil { // has else
		return false
	}

	cond, ok := v.Cond.(*ast.Ident)
	if !ok { // no binary expression
		return false
	}

	return slices.Contains(names, cond.Name)
}
