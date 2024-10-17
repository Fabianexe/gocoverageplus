package entity

import (
	"fmt"
	"strconv"
)

type LineCounter struct {
	totalLines   int
	coveredLines int
}

func (c *LineCounter) AddLine(covered bool) {
	c.totalLines++
	if covered {
		c.coveredLines++
	}
}

func (c *LineCounter) String() string {
	if c.totalLines == 0 {
		return "1.00"
	}

	return fmt.Sprintf("%.2f", float64(c.coveredLines)/float64(c.totalLines))
}

func (c *LineCounter) ValidString() string {
	return strconv.Itoa(c.totalLines)
}

func (c *LineCounter) CoveredString() string {
	return strconv.Itoa(c.coveredLines)
}

type BranchCounter struct {
	totalBranches   int
	coveredBranches int
}

func (b *BranchCounter) AddBranch(covered bool) {
	b.totalBranches++
	if covered {
		b.coveredBranches++
	}
}

func (b *BranchCounter) String() string {
	if b.totalBranches == 0 {
		return "1.00"
	}

	return fmt.Sprintf("%.2f", float64(b.coveredBranches)/float64(b.totalBranches))
}

func (b *BranchCounter) ValidString() string {
	return strconv.Itoa(b.totalBranches)
}

func (b *BranchCounter) CoveredString() string {
	return strconv.Itoa(b.coveredBranches)
}
