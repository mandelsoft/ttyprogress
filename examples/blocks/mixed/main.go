package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress/blocks"
)

func Download(blocks *blocks.Blocks) {
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

func Process(blocks *blocks.Blocks, id int) {
	total := 10 + rand.Int()%20
	w := blocks.NewBlock(3)
	go func() {
		done := 0
		for done < total {
			done++
			fmt.Fprintf(w, "reached state %d:%d\n", id, done)
			time.Sleep(time.Millisecond * time.Duration(100+rand.Int()%500))
		}
		w.Close()
	}()
}

func main() {
	blocks := blocks.New(os.Stdout)

	for b := 0; b < 5; b++ {
		if rand.Int()%2 == 0 {
			Download(blocks)
		} else {
			Process(blocks, b)
		}
	}
	blocks.Close()
	blocks.Wait(nil)
}
