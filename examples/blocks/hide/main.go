package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress/blocks"
)

func main() {
	blks := blocks.New(os.Stdout)

	var list [3]*blocks.Block

	for b := 0; b < len(list); b++ {
		list[b] = blks.NewBlock(1).HideOnClose()
		list[b].Reset()
		fmt.Fprintf(list[b], "Downloading %d\n", b+1)
		list[b].Flush()
		time.Sleep(time.Second)
	}
	blks.Close()
	time.Sleep(3 * time.Second)
	list[1].Close()
	time.Sleep(3 * time.Second)
	list[0].Close()
	time.Sleep(3 * time.Second)
	list[2].Close()
}
