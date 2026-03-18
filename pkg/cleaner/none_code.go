package cleaner

import (
	"go/ast"
	"go/token"
	"log/slog"

	"github.com/Fabianexe/gocoverageplus/pkg/entity"
)

func cleanNoneCodeLines(project *entity.Project) *entity.Project {
	slog.Info("Clean none code lines")
	for _, p := range project.Packages {
		for _, f := range p.Files {
			var cleanedLines int
			for _, method := range f.Methods {
				noneCodeVisitor := &noneCodeVisitor{
					validLines: make(map[int]struct{}, 128),
					fset:       p.Fset,
				}

				ast.Walk(noneCodeVisitor, method.Body)

				for line := method.Tree.StartPosition.Line; line < method.Tree.EndPosition.Line; line++ {
					if _, ok := noneCodeVisitor.validLines[line]; !ok {
						cleanedLines++
						start := method.File.Position(method.File.LineStart(line))
						end := method.File.Position(method.File.LineStart(line+1) - 1)
						method.Tree.AddBlock(
							&entity.Block{
								StartPosition: start,
								EndPosition:   end,
								DefPosition:   start,
								Type:          entity.TypeBlock,
								Ignore:        true,
							},
						)
					}
				}
			}
			if cleanedLines > 0 {
				slog.Debug("Cleaned lines", "File", f.FilePath, "Lines", cleanedLines)
			}
		}
	}

	return project
}

type noneCodeVisitor struct {
	validLines map[int]struct{}
	fset       *token.FileSet
}

func (n *noneCodeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		return n
	}

	lineNUmber := n.fset.Position(node.Pos()).Line
	n.validLines[lineNUmber] = struct{}{}

	return n
}
