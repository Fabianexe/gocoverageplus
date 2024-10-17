// Package entity contains all entities that are used in every component of the application
package entity

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/cover"
)

type Project struct {
	Packages       []*Package
	LineCoverage   LineCounter
	BranchCoverage BranchCounter
}

type Package struct {
	Name           string
	Files          []*File
	Fset           *token.FileSet
	LineCoverage   LineCounter
	BranchCoverage BranchCounter
}

type File struct {
	Name           string
	FilePath       string
	Ast            *ast.File
	Methods        []*Method
	LineCoverage   LineCounter
	BranchCoverage BranchCounter
}

type Method struct {
	Name           string
	Body           *ast.BlockStmt
	Tree           *Block
	LineCoverage   LineCounter
	BranchCoverage BranchCounter
	Complexity     int
	File           *token.File
	cover          []cover.ProfileBlock
	branches       []*Branch
	lines          []*Line
}

type Line struct {
	Number        int
	CoverageCount int
}

type Branch struct {
	DefLine int
	Covered bool
}

func (m *Method) GetBranches() []*Branch {
	if m.branches == nil {
		m.branches = m.getBranches(m.Tree)
	}

	return m.branches
}

func (m *Method) getBranches(b *Block) []*Branch {
	branches := make([]*Branch, 0, 128)
	if b.Type == TypeBranch {
		covered := false
		if len(b.Coverage) > 0 {
			covered = true
		}
		branches = append(branches, &Branch{
			DefLine: b.DefPosition.Line,
			Covered: covered,
		})
	}

	for _, child := range b.Children {
		branches = append(branches, m.getBranches(child)...)
	}

	return branches
}

func (m *Method) GetLines() []*Line {
	if m.lines == nil {
		m.lines = m.getLines()
	}

	return m.lines
}

func (m *Method) getLines() []*Line {
	lines := make([]*Line, 0, m.Tree.EndPosition.Line-m.Tree.DefPosition.Line+1)
	covers := m.GetCover()
	dict := make(map[int]*Line, cap(lines))
	for _, c := range covers {
		for i := c.StartLine; i <= c.EndLine; i++ {
			if existing, ok := dict[i]; !ok {
				l := &Line{
					Number:        i,
					CoverageCount: c.Count,
				}
				dict[i] = l
				lines = append(lines, l)
			} else {
				existing.CoverageCount = max(existing.CoverageCount, c.Count)
			}
		}
	}

	return lines
}
