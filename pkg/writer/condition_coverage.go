package writer

import (
	"fmt"
)

type conditionCoverage struct {
	total   int
	covered int
}

func (c *conditionCoverage) Add(covered bool) {
	c.total++

	if covered {
		c.covered++
	}
}

func (c *conditionCoverage) String() string {
	return fmt.Sprintf("%.2f (%d/%d)", float64(c.covered)/float64(c.total), c.covered, c.total)
}
