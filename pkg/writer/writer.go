// Package writer writes xml file based on the cobertura dtd:
// https://github.com/cobertura/cobertura/blob/master/cobertura/src/test/resources/dtds/coverage-04.dtd
package writer

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/Fabianexe/gocoverageplus/pkg/entity"
)

func WriteXML(path string, project *entity.Project, outPath string) error {
	xmlCoverage := ConvertToCobertura(path, project)

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}

	encoder := xml.NewEncoder(outFile)
	encoder.Indent("", "\t")

	slog.Info("Write coverage to file", "Path", outPath)
	err = encoder.Encode(xmlCoverage)
	if err != nil {
		return err
	}

	if err := outFile.Close(); err != nil {
		return err
	}

	return nil
}

func WriteTextFMT(project *entity.Project, outPath string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}

	slog.Info("Write coverage to file", "Path", outPath)
	if _, err := outFile.WriteString("mode: set\n"); err != nil {
		return err
	}
	for _, pkg := range project.Packages {
		for _, file := range pkg.Files {
			for _, method := range file.Methods {
				for _, block := range method.GetCover() {
					if _, err := outFile.WriteString(
						fmt.Sprintf(
							"%s/%s:%d.%d,%d.%d %d %d\n",
							pkg.Name,
							path.Base(file.FilePath),
							block.StartLine,
							block.StartCol,
							block.EndLine,
							block.EndCol,
							block.NumStmt,
							block.Count,
						),
					); err != nil {
						return err
					}
				}

			}
		}
	}

	if err := outFile.Close(); err != nil {
		return err
	}

	return nil
}
