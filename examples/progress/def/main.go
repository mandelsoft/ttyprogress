package main

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
)

var BarElapsed = ttyprogress.TypeFor(ttyprogress.NewBar(5).SetWidth(50).AppendElapsed())

func run(p ttyprogress.Context, bar *ttyprogress.BarDefinition) error {
	n := 10 + rand.Int()%10
	e, err := bar.SetTotal(n).Add(p)
	if err != nil {
		return err
	}
	go func() {
		e.Start()
		for i := 0; i < n; i++ {
			time.Sleep(time.Second + time.Millisecond*100*time.Duration(rand.Intn(10)))
			e.Incr()
		}
		e.Close()
	}()
	return nil
}

func main() {
	p := ttyprogress.For(os.Stdout)

	bar1 := ttyprogress.New(BarElapsed).PrependFunc(ttyprogress.Message("first action "))
	bar2 := ttyprogress.New(BarElapsed).PrependFunc(ttyprogress.Message("second action"))

	run(p, bar1)
	run(p, bar2)
	p.Close()

	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	p.Wait(ctx)

}
