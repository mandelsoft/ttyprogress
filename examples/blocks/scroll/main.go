package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress/blocks"
)

func main() {
	blocks := blocks.New(os.Stdout)

	writer := blocks.NewBlock(3).SetAuto().SetGap("-> ").SetTitleLine("Some work:")

	for i := 0; i <= 20; i++ {
		fmt.Fprintf(writer, "doing step %d\n", i)
		time.Sleep(time.Millisecond * 500)
	}

	writer.Close()
}
