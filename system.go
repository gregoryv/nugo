package graph

import "github.com/gregoryv/graph/internal"

func NewSystem() *System {
	rn := internal.NewRoot("/")
	rn.Make("/bin")
	return &System{
		rn: rn,
	}
}

type System struct {
	rn *internal.RootNode
}
