package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress/blocks"
)

func main() {
	blocks := blocks.New(os.Stdout)

	writer := blocks.NewBlock(2).SetFinal("Finished: Downloaded 100GB")

	for i := 0; i <= 20; i++ {
		writer.Reset()
		fmt.Fprintf(writer, "Downloading.. (%d/%d) GB\n", i*5, 100)
		writer.Flush()
		time.Sleep(time.Millisecond * 500)
	}

	writer.Close()
}
