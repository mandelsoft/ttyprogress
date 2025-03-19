package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress/blocks"
)

func main() {
	blocks := blocks.New(os.Stdout)

	for b := 0; b < 3; b++ {
		total := 100 + rand.Int()%100
		w := blocks.NewBlock(1).SetFinal(fmt.Sprintf("Finished: Downloaded %d GB", total))

		go func() {
			done := 0
			for done < total {
				w.Reset()
				fmt.Fprintf(w, "Downloading.. (%d/%d) GB\n", done, total)
				w.Flush()
				time.Sleep(time.Millisecond * time.Duration(100+rand.Int()%500))
				done = done + 5
			}
			w.Close()
		}()
	}
	blocks.Close()
	blocks.Wait(nil)
}
