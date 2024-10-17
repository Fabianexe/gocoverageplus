package complexity

import (
	"go/ast"
)

func getCyclomaticComplexity(root ast.Node) int {
	visitor := &cyclomaticVisitor{
		complexity: 1,
	}

	ast.Walk(visitor, root)

	return visitor.complexity

}

type cyclomaticVisitor struct {
	complexity int
}

func (c *cyclomaticVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch node.(type) {
	case *ast.IfStmt,
		*ast.ForStmt,
		*ast.RangeStmt,
		*ast.FuncDecl,
		*ast.SwitchStmt,
		*ast.TypeSwitchStmt,
		*ast.SelectStmt:
		c.complexity++
	}

	return c
}
