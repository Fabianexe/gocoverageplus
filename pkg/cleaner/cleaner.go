// Package cleaner cleans the coverage data and discard lines, classes and packages that are not relevant
package cleaner

import (
	"github.com/Fabianexe/gocoverageplus/pkg/entity"
)

// CleanData cleans the package data
func CleanData(
	project *entity.Project,
	cGeneratedFiles bool,
	cNoneCodeLines bool,
	cErrorIf bool,
	customIfNames []string,
) *entity.Project {
	if cGeneratedFiles {
		project = cleanGeneratedFiles(project)
	}

	if cNoneCodeLines {
		project = cleanNoneCodeLines(project)
	}

	if cErrorIf {
		project = cleanErrorIf(project)
	}
	if len(customIfNames) > 0 {
		project = cleanCustomIf(project, customIfNames)
	}

	return project
}
