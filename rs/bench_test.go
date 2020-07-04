package rs

import (
	"fmt"
	"testing"
)

func BenchmarkSyscall_mkdir(b *testing.B) {
	var (
		sys = &Syscall{System: NewSystem(), acc: Root}
	)
	for i := 0; i < b.N; i++ {
		sys.mkdir(fmt.Sprintf("/dir%d", i), 0)
	}
}
