// Package source reads source files from a given path and create entities from this
package source

import (
	"go/ast"
	"log/slog"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"

	"github.com/Fabianexe/gocoverageplus/pkg/entity"
)

func LoadSources(path string) (*entity.Project, error) {
	goPackages, err := getGoPaths(path)
	if err != nil {
		return nil, err
	}

	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedModule |
			packages.NeedTypes,
		Dir: path,
	}

	pkgs, err := packages.Load(cfg, goPackages...)
	if err != nil {
		return nil, err
	}

	var countPackages, countFiles, countMethods int
	allPackages := make([]*entity.Package, 0, len(pkgs))
	for _, pkg := range pkgs {
		pack := &entity.Package{
			Name:  pkg.PkgPath,
			Files: make([]*entity.File, 0, len(pkg.Syntax)),
			Fset:  pkg.Fset,
		}

		slog.Debug("Package", "Path", pkg.PkgPath, "Files", len(pkg.Syntax))
		for i, fileAst := range pkg.Syntax {
			methodsMap := make(map[string][]*entity.Method)
			for _, decl := range fileAst.Decls {
				if fun, ok := decl.(*ast.FuncDecl); ok {
					method := &entity.Method{
						Name: fun.Name.Name,
						Body: fun.Body,
						File: pkg.Fset.File(fun.Pos()),
					}

					// start after the function declaration
					startLine := pkg.Fset.Position(fun.Body.Lbrace).Line + 1
					endLine := pkg.Fset.Position(fun.End()).Line
					if startLine >= endLine {
						continue
					}

					bV := &branchVisitor{
						fset: pkg.Fset,
					}

					ast.Walk(bV, fun)

					method.Tree = bV.getTree()

					countMethods++
					className := getClassName(fun)
					methodsMap[className] = append(methodsMap[className], method)
				}
			}

			var methodCount int
			for className, methods := range methodsMap {
				file := &entity.File{
					Name:     className,
					FilePath: pkg.GoFiles[i],
					Ast:      fileAst,
					Methods:  methods,
				}
				pack.Files = append(pack.Files, file)

				methodCount += len(methods)
			}

			slog.Debug("File", "Name", filepath.Base(pkg.GoFiles[i]), "Methods", methodCount)

			countFiles++
		}

		countPackages++
		allPackages = append(allPackages, pack)
	}
	slog.Info("Source reading Finished", "Packages", countPackages, " Files", countFiles, " Methods", countMethods)
	return &entity.Project{Packages: allPackages}, nil
}

func getGoPaths(path string) ([]string, error) {
	goPath := make(map[string]struct{}, 1000)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" {
			goPath[filepath.Dir(path)] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	goPackages := make([]string, 0, len(goPath))
	for pack := range goPath {
		goPackages = append(goPackages, pack)
	}
	return goPackages, nil
}

func getClassName(fun *ast.FuncDecl) string {
	if fun.Recv == nil {
		return "-"
	}

	if star, ok := fun.Recv.List[0].Type.(*ast.StarExpr); ok {
		if index, ok := star.X.(*ast.IndexExpr); ok {
			return index.X.(*ast.Ident).Name
		}

		if index, ok := star.X.(*ast.IndexListExpr); ok {
			return index.X.(*ast.Ident).Name
		}

		return star.X.(*ast.Ident).Name
	}

	if index, ok := fun.Recv.List[0].Type.(*ast.IndexExpr); ok {
		return index.X.(*ast.Ident).Name
	}

	if index, ok := fun.Recv.List[0].Type.(*ast.IndexListExpr); ok {
		return index.X.(*ast.Ident).Name
	}

	return fun.Recv.List[0].Type.(*ast.Ident).Name
}
