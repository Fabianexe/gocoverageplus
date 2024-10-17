package complexity

import (
	"go/ast"
)

func getCognitiveComplexity(root ast.Node) int {
	visitor := &cognitiveVisitor{
		complexity: 1,
	}

	ast.Walk(visitor, root)

	return visitor.complexity

}

type cognitiveVisitor struct {
	complexity int
	level      []ast.Node
}

func (c *cognitiveVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		return c
	}

	for len(c.level) > 0 && c.level[len(c.level)-1].End() < node.Pos() {
		c.level = c.level[:len(c.level)-1]
	}
	switch node.(type) {
	case *ast.IfStmt,
		*ast.ForStmt,
		*ast.RangeStmt,
		*ast.FuncDecl,
		*ast.SwitchStmt,
		*ast.TypeSwitchStmt,
		*ast.SelectStmt:
		c.complexity += len(c.level) + 1 // function is the missing level
		c.level = append(c.level, node)
	}

	return c
}
