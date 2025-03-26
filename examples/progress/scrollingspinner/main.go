package main

import (
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
)

func main() {
	p := ttyprogress.For(os.Stdout)

	bar, _ := ttyprogress.NewScrollingSpinner("doing some calculations", 10).
		SetDone("calculations done").
		AppendElapsed().
		Add(p)
	bar.Start()
	p.Close()

	time.Sleep(time.Millisecond * 11000)
	bar.Close()
}
