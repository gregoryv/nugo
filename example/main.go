package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/gregoryv/nugo"
)

func main() {
	cf, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cf)
	defer pprof.StopCPUProfile()

	// to profile
	_, rn := largeTree()
	stop := time.After(5 * time.Second)
loop:
	for {
		select {
		case <-stop:
			break loop
		default:
			rn.Find("/9/9/9")
		}
	}

	// Save profile
	f, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(f)
	f.Close()
}

func largeTree() (int, *nugo.Node) {
	rn := nugo.NewRootNode("/", nugo.ModeSort|nugo.ModeDistinct)
	var count int
	addRec(rn, 10, 4, &count)
	return count, rn
}

func addRec(parent *nugo.Node, nodes, level int, count *int) {
	if level == 0 {
		return
	}
	for i := 0; i < nodes; i++ {
		child := nugo.NewNode(fmt.Sprintf("%v", i))
		parent.Add(child)
		*count = *count + 1
		addRec(child, nodes, (level - 1), count)
	}
}
