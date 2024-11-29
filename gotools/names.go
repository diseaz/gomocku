package gotools

import "fmt"

type NameGenerator struct {
	prefix string
	counts map[string]int
}

func NewNameGenerator(prefix string) *NameGenerator {
	return &NameGenerator{
		prefix: prefix,
		counts: make(map[string]int),
	}
}

func (g *NameGenerator) Next() string {
	return g.NextForPrefix(g.prefix)
}

func (g *NameGenerator) NextForPrefix(p string) string {
	g.counts[p]++
	return fmt.Sprintf("%s%d", p, g.counts[p])

}
