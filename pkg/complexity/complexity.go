// Package complexity enriches the enties with complexity metrics
package complexity

import (
	"log/slog"

	"github.com/Fabianexe/gocoverageplus/pkg/entity"
)

// AddComplexity adds complexity metrics to the packages
func AddComplexity(project *entity.Project, useCyclomaticComplexity bool) *entity.Project {
	if useCyclomaticComplexity {
		slog.Info("Use cyclomatic complexity")
	} else {
		slog.Info("Use cognitive complexity")
	}
	for _, p := range project.Packages {
		for _, f := range p.Files {
			for _, method := range f.Methods {
				if useCyclomaticComplexity {
					method.Complexity = getCyclomaticComplexity(method.Body)
				} else {
					method.Complexity = getCognitiveComplexity(method.Body)
				}
			}
		}
	}

	return project
}
