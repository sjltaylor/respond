package testHelpers

import (
	"fmt"
)

type StringGenerator struct {
	fmt string
	seq int64
}

func NewStringGenerator (fmt string) *StringGenerator {
	return &StringGenerator{
		fmt: fmt,
		seq: 1,
	}
}

func (g *StringGenerator) Next() string {
	g.seq++
	return g.With(g.seq)
}

func (g *StringGenerator) With(i int64) string {
	return fmt.Sprintf(g.fmt, i)
}
