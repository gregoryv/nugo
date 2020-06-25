package graph

func NewSystem() *System {
	rn := NewRoot("/")
	rn.Make("/bin")
	return &System{
		rn: rn,
	}
}

type System struct {
	rn *rootNode
}
