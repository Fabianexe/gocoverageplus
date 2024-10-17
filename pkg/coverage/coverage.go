// Package coverage loads a golang coverage report and enrich the entities with the information
package coverage

import (
	"log/slog"
	"path/filepath"
	"strings"

	"golang.org/x/tools/cover"

	"github.com/Fabianexe/gocoverageplus/pkg/entity"
)

// LoadCoverage loads the coverage data from the given file
func LoadCoverage(project *entity.Project, coverageReport string) (*entity.Project, error) {
	profiles, err := cover.ParseProfiles(coverageReport)
	if err != nil {
		return nil, err
	}

	for _, p := range profiles {
		slog.Debug("Profile", "Path", p.FileName, "Blocks", len(p.Blocks))
		found := false
		for _, pack := range project.Packages {
			if !strings.HasPrefix(p.FileName, pack.Name) {
				continue
			}
			found = true
			filename := filepath.Base(p.FileName)
			for _, f := range pack.Files {
				if filepath.Base(f.FilePath) != filename {
					continue
				}
				applyBlocks(f.Methods, p.Blocks)
			}

		}
		if !found {
			slog.Warn("Not found source for: " + p.FileName)
		}
	}

	updateLineCoverage(project)
	updateBranchCoverage(project)

	return project, nil
}

func applyBlocks(methods []*entity.Method, blocks []cover.ProfileBlock) {
	for _, b := range blocks {
		if b.Count == 0 {
			continue
		}
		for _, method := range methods {
			method.Tree.AddCoverageBlock(b)
		}
	}
}

func updateLineCoverage(project *entity.Project) {
	for _, pack := range project.Packages {
		for _, f := range pack.Files {
			for _, method := range f.Methods {
				for _, line := range method.GetLines() {
					isCovered := line.CoverageCount > 0
					method.LineCoverage.AddLine(isCovered)
					f.LineCoverage.AddLine(isCovered)
					pack.LineCoverage.AddLine(isCovered)
					project.LineCoverage.AddLine(isCovered)

				}
			}
		}
	}
}

func updateBranchCoverage(project *entity.Project) {
	for _, pack := range project.Packages {
		for _, f := range pack.Files {
			for _, method := range f.Methods {
				for _, branch := range method.GetBranches() {
					isCovered := branch.Covered
					method.BranchCoverage.AddBranch(isCovered)
					f.BranchCoverage.AddBranch(isCovered)
					pack.BranchCoverage.AddBranch(isCovered)
					project.BranchCoverage.AddBranch(isCovered)

				}
			}
		}
	}
}
