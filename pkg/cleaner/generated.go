package cleaner

import (
	"go/ast"
	"log/slog"

	"github.com/Fabianexe/gocoverageplus/pkg/entity"
)

func cleanGeneratedFiles(project *entity.Project) *entity.Project {
	slog.Info("Clean generated files")
	for _, p := range project.Packages {
		i := 0
		for i < len(p.Files) {
			f := p.Files[i]
			if ast.IsGenerated(f.Ast) {
				slog.Debug("Remove generated file", "File", f.FilePath)
				p.Files = append(p.Files[:i], p.Files[i+1:]...)
				continue
			}
			i++
		}
	}

	return project
}
