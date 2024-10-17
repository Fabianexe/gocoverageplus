package entity

import (
	"go/ast"
	"go/token"
	"slices"

	"golang.org/x/tools/cover"
)

func (m *Method) GetCover() []cover.ProfileBlock {
	if m.cover != nil {
		return m.cover
	}
	coverBlocks := m.generateCover()
	vis := &posGatherVisitor{
		lines: make(map[int][]token.Pos),
		file:  m.File,
	}
	ast.Walk(vis, m.Body)
	minLines, maxLines := vis.getLines()

	retBlocks := make([]cover.ProfileBlock, 0, len(coverBlocks))
	for _, block := range coverBlocks {
		if block.NumStmt > 0 {
			maxLine, ok := maxLines[block.StartLine]
			if !ok || maxLine < block.StartCol {
				pos := m.File.Position(m.File.LineStart(block.StartLine + 1))
				block.StartLine = pos.Line
				block.StartCol = pos.Column
			}
			minLine, ok := minLines[block.EndLine]
			if !ok || minLine > block.EndCol {
				pos := m.File.Position(m.File.LineStart(block.EndLine) - 1)
				block.EndLine = pos.Line
				block.EndCol = pos.Column
			}

			retBlocks = append(retBlocks, block)
		}
	}

	m.cover = retBlocks

	return m.cover
}

func (m *Method) generateCover() []cover.ProfileBlock {
	coverBlocks, _, _, _, _ := m.internalCover(m.Tree)

	return coverBlocks
}

func (m *Method) internalCover(b *Block) (results []cover.ProfileBlock, cutFrom, cutTo token.Position, statmentsBefore, statmentsAfter int) {
	if b.Ignore {
		cutFrom = b.DefPosition
		cutTo = b.EndPosition

		return
	}
	switch b.Type {
	case TypeBlock:
		return m.generateCoverBlock(b)
	case TypeBranch:
		return m.generateCoverBlock(b)
	case TypeAtomic:
		return m.generateCoverAtomic(b)
	}

	panic("Unknown block type")
}

func (m *Method) generateCoverBlock(b *Block) (results []cover.ProfileBlock, cutFrom, cutTo token.Position, statmentsBefore, statmentsAfter int) {
	cutFrom = b.StartPosition
	cutTo = b.EndPosition
	results = make([]cover.ProfileBlock, 0, 128)
	lastEndPos := m.movePos(b.StartPosition, 1)
	for _, child := range b.Children {
		childResults, childCutFrom, childCutTo, childStatmentsBefore, childStatmentsAfter := m.internalCover(child)
		// valid cut
		if childCutFrom.Line > 0 {
			if childCutFrom.Offset > lastEndPos.Offset {
				results = append(
					results,
					b.createProfileBlock(lastEndPos, childCutFrom, statmentsAfter+childStatmentsBefore),
				)
			}
			lastEndPos = childCutTo
			statmentsAfter = 0
		}
		statmentsAfter += childStatmentsAfter
		results = append(results, childResults...)
	}

	if lastEndPos.Offset < b.EndPosition.Offset {
		results = append(
			results,
			b.createProfileBlock(lastEndPos, m.movePos(b.EndPosition, -1), statmentsAfter),
		)
		statmentsAfter = 0
	}

	return
}

func (m *Method) generateCoverAtomic(b *Block) (results []cover.ProfileBlock, cutFrom, cutTo token.Position, statmentsBefore, statmentsAfter int) {
	if len(b.Children) == 0 {
		statmentsAfter = 1

		return
	}

	cutTo = b.EndPosition
	results = make([]cover.ProfileBlock, 0, 128)
	isFirst := true
	for _, child := range b.Children {
		childResults, childCutFrom, childCutTo, childStatmentsBefore, childStatmentsAfter := m.internalCover(child)
		// valid cut
		if childCutFrom.Line > 0 {
			if isFirst {
				statmentsBefore = statmentsAfter + childStatmentsBefore
				isFirst = false
				cutFrom = childCutFrom
			} else {
				results = append(
					results,
					b.createProfileBlock(cutTo, childCutFrom, statmentsAfter+childStatmentsBefore),
				)
			}

			cutTo = childCutTo
			statmentsAfter = 0
		}
		statmentsAfter += childStatmentsAfter
		results = append(results, childResults...)
	}

	return
}

func (m *Method) movePos(pos token.Position, move int) token.Position {
	p := m.File.Pos(pos.Offset + move)

	return m.File.Position(p)
}

type posGatherVisitor struct {
	lines map[int][]token.Pos
	file  *token.File
}

func (n *posGatherVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		return n
	}

	lineNUmber := n.file.Position(node.Pos()).Line
	n.lines[lineNUmber] = append(n.lines[lineNUmber], node.Pos())

	return n
}

func (n *posGatherVisitor) getLines() (minLines, maxLines map[int]int) {
	minLines = make(map[int]int, len(n.lines))
	maxLines = make(map[int]int, len(n.lines))
	for lines, pos := range n.lines {
		minPos := slices.Min(pos)
		maxPos := slices.Max(pos)
		minLines[lines] = n.file.Position(minPos).Column
		maxLines[lines] = n.file.Position(maxPos).Column
	}

	return
}
