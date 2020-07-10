package nugo

import (
	"testing"
)

var (
	rn    *Node
	nodes int
)

func init() {
	rn = NewRootNode("/", ModeSort|ModeDistinct)
	for _, name := range "abcdefghihklmn" {
		nodes++
		c := NewNode(string(name))
		rn.Add(c)
		for _, name := range "1234567890" {
			nodes++
			c.Add(NewNode(string(name)))
		}
	}
}

func BenchmarkRootNode_Find_first(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rn.Find("/")
	}
}

func BenchmarkRootNode_Find_last(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rn.Find("/n/0")
	}
}
