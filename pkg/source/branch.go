package source

import (
	"go/ast"
	"go/token"

	"github.com/Fabianexe/gocoverageplus/pkg/entity"
)

type branchVisitor struct {
	blocks []*entity.Block
	fset   *token.FileSet
}

func (b *branchVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch v := node.(type) {
	// block
	case *ast.FuncDecl:
		b.blocks = append(b.blocks, b.createBranch(v.Pos(), v.Body.Pos(), v.Body.End(), entity.TypeBlock))
	case *ast.FuncLit:
		b.blocks = append(b.blocks, b.createBranch(v.Pos(), v.Pos(), v.Body.Pos(), entity.TypeAtomic))
		b.blocks = append(b.blocks, b.createBranch(v.Pos(), v.Body.Pos(), v.Body.End(), entity.TypeBlock))
	case *ast.BlockStmt:
		b.blocks = append(b.blocks, b.createBranch(v.Pos(), v.Pos(), v.End(), entity.TypeBlock))
	// branch
	case *ast.CaseClause:
		b.blocks = append(b.blocks, b.createBranch(v.Pos(), v.Colon, v.End(), entity.TypeBranch))
	case *ast.CommClause:
		b.blocks = append(b.blocks, b.createBranch(v.Pos(), v.Colon, v.End(), entity.TypeBranch))
	case *ast.IfStmt:
		b.blocks = append(b.blocks, b.createBranch(v.Pos(), v.Pos(), v.Body.Pos(), entity.TypeAtomic))
		b.blocks = append(b.blocks, b.createBranch(v.Pos(), v.Body.Pos(), v.Body.End(), entity.TypeBranch))
	case *ast.ForStmt:
		b.blocks = append(b.blocks, b.createBranch(v.Pos(), v.Pos(), v.Body.Pos(), entity.TypeAtomic))
		b.blocks = append(b.blocks, b.createBranch(v.Pos(), v.Body.Pos(), v.Body.End(), entity.TypeBranch))
	case *ast.RangeStmt:
		b.blocks = append(b.blocks, b.createBranch(v.Pos(), v.Pos(), v.Body.Pos(), entity.TypeAtomic))
		b.blocks = append(b.blocks, b.createBranch(v.Pos(), v.Body.Pos(), v.Body.End(), entity.TypeBranch))
	// atomic
	default:
		if v == nil {
			return b
		}

		b.blocks = append(b.blocks, b.createBranch(v.Pos(), v.Pos(), v.End(), entity.TypeAtomic))
	}

	return b
}

func (b *branchVisitor) createBranch(def, start, end token.Pos, t entity.BlockType) *entity.Block {
	return &entity.Block{
		DefPosition:   b.fset.Position(def),
		StartPosition: b.fset.Position(start),
		EndPosition:   b.fset.Position(end),
		Type:          t,
	}
}

func (b *branchVisitor) getBlocks() []*entity.Block {
	retBlocks := make([]*entity.Block, 0, len(b.blocks))
	last := &entity.Block{}
	for _, block := range b.blocks {
		if block.DefPosition.Line == last.DefPosition.Line && block.EndPosition.Line == last.EndPosition.Line {
			if last.Type < block.Type {
				retBlocks[len(retBlocks)-1] = block
				last = block
			}
			continue
		}
		retBlocks = append(retBlocks, block)
		last = block
	}

	return retBlocks
}

func (b *branchVisitor) getTree() *entity.Block {
	blocks := b.getBlocks()
	rootB := blocks[0]
	root := &entity.Block{
		StartPosition: rootB.StartPosition,
		EndPosition:   rootB.EndPosition,
		DefPosition:   rootB.DefPosition,
		Type:          rootB.Type,
	}

	last := root
	for _, block := range blocks[1:] {
		for block.EndPosition.Line > last.EndPosition.Line {
			last = last.Parent
		}
		newBlock := &entity.Block{
			StartPosition: block.StartPosition,
			EndPosition:   block.EndPosition,
			DefPosition:   block.DefPosition,
			Type:          block.Type,
			Parent:        last,
		}
		last.Children = append(last.Children, newBlock)
		last = newBlock
	}

	return root
}
