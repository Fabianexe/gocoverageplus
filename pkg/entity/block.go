package entity

import (
	"go/token"

	"golang.org/x/tools/cover"
)

type BlockType uint8

const (
	TypeAtomic BlockType = iota
	TypeBlock
	TypeBranch
)

type Block struct {
	StartPosition token.Position
	EndPosition   token.Position
	DefPosition   token.Position
	Coverage      []cover.ProfileBlock
	Type          BlockType
	Children      []*Block
	Parent        *Block
	Ignore        bool
}

func (b *Block) AddCoverageBlock(block cover.ProfileBlock) {
	if block.EndLine < b.StartPosition.Line ||
		block.EndLine == b.StartPosition.Line && block.EndCol <= b.StartPosition.Column ||
		block.StartLine > b.EndPosition.Line ||
		block.StartLine == b.EndPosition.Line && block.StartCol >= b.EndPosition.Column {
		return
	}
	if b.Type == TypeAtomic {
		if len(b.Coverage) == 0 {
			b.Coverage = append(b.Coverage, cover.ProfileBlock{
				StartLine: b.StartPosition.Line,
				StartCol:  b.StartPosition.Column,
				EndLine:   b.EndPosition.Line,
				EndCol:    b.EndPosition.Column,
			})
		}
		b.Coverage[0].Count = max(b.Coverage[0].Count, block.Count)
	} else {
		b.Coverage = append(b.Coverage, block)
	}
	for _, child := range b.Children {
		child.AddCoverageBlock(block)
	}
}

func (b *Block) AddBlock(newB *Block) {
	if len(b.Children) == 0 {
		b.Children = []*Block{newB}

		return
	}

	newChilds := make([]*Block, 0, len(b.Children)+1)
	added := false
	for _, child := range b.Children {
		// child is part of new block so ignore it
		if child.StartPosition.Offset >= newB.StartPosition.Offset &&
			child.EndPosition.Offset <= newB.EndPosition.Offset {
			continue
		}
		// new Block is part of child so add it there:
		if !added &&
			child.StartPosition.Offset < newB.StartPosition.Offset &&
			child.EndPosition.Offset > newB.EndPosition.Offset {
			child.AddBlock(newB)
			added = true
		}
		// new Block is not added and before current child
		if !added &&
			child.StartPosition.Offset > newB.EndPosition.Offset {
			newB.Parent = b
			newChilds = append(newChilds, newB)
			added = true
		}

		newChilds = append(newChilds, child)
	}
	if !added {
		newB.Parent = b
		newChilds = append(newChilds, newB)
	}

	b.Children = newChilds
}

func (b *Block) createProfileBlock(
	startPos token.Position,
	endPos token.Position,
	statements int,
) cover.ProfileBlock {
	cov := 0
	for _, block := range b.Coverage {
		if block.EndLine < startPos.Line ||
			block.EndLine == startPos.Line && block.EndCol <= startPos.Column ||
			block.StartLine > endPos.Line ||
			block.StartLine == endPos.Line && block.StartCol >= endPos.Column {
			continue
		}

		cov = max(cov, block.Count)
	}

	return cover.ProfileBlock{
		StartLine: startPos.Line,
		StartCol:  startPos.Column,
		EndLine:   endPos.Line,
		EndCol:    endPos.Column,
		NumStmt:   statements,
		Count:     cov,
	}
}
