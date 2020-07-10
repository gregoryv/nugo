package nugo

import "strings"

// A NodeMode represents a node's mode and permission bits.
type NodeMode uint32

const (
	ModeDir      NodeMode = 1 << (32 - 1 - iota)
	ModeSort              // sorted by name
	ModeDistinct          // no duplicate children
	ModeRoot

	ModeType NodeMode = ModeSort | ModeDistinct
	ModePerm NodeMode = 07777
)

// String returns the bits as string using rwx notation for each bit.
func (m NodeMode) String() string {
	var s strings.Builder
	s.WriteByte(mChar(m, ModeDir, 'd'))
	s.WriteByte(mChar(m, 04000, 'r'))
	s.WriteByte(mChar(m, 02000, 'w'))
	s.WriteByte(mChar(m, 01000, 'x'))
	s.WriteByte(mChar(m, 00400, 'r'))
	s.WriteByte(mChar(m, 00200, 'w'))
	s.WriteByte(mChar(m, 00100, 'x'))
	s.WriteByte(mChar(m, 00040, 'r'))
	s.WriteByte(mChar(m, 00020, 'w'))
	s.WriteByte(mChar(m, 00010, 'x'))
	s.WriteByte(mChar(m, 00004, 'r'))
	s.WriteByte(mChar(m, 00002, 'w'))
	s.WriteByte(mChar(m, 00001, 'x'))
	return s.String()
}

func mChar(m, mask NodeMode, c byte) byte {
	if m&mask == mask {
		return c
	}
	return '-'
}
