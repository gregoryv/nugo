package sys

import (
	"fmt"
	"testing"
)

func BenchmarkSyscall_Mkdir(b *testing.B) {
	var (
		sys = &Syscall{System: NewSystem(), acc: Root}
	)
	for i := 0; i < b.N; i++ {
		sys.Mkdir(fmt.Sprintf("/dir%d", i), 0)
	}
}
