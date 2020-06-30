package rs

func Example_setupNewSystem() {
	var (
		sys    = NewSystem()
		asRoot = sys.Use(Root)
	)
	asRoot.ExecCmd("/bin/mkdir", "/tmp")
	// output:
}
