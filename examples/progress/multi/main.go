package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mandelsoft/ttyprogress"
	"github.com/mandelsoft/ttyprogress/units"
)

func main() {
	p := ttyprogress.For(os.Stdout)

	for b := 0; b < 3; b++ {
		total := 100 + rand.Int()%100
		w, _ := ttyprogress.NewBar().
			SetTotal(total).
			SetPredefined(1).
			SetFinal(fmt.Sprintf("Finished: Downloaded %d GB", total)).
			AppendFunc(ttyprogress.Amount(units.Bytes(units.GB))).
			PrependFunc(ttyprogress.Message("Downloading ...")).
			Add(p)

		go func() {
			done := 0
			for w.Set(done) {
				time.Sleep(time.Millisecond * time.Duration(100+rand.Int()%500))
				done = done + 5
			}
		}()
	}
	p.Close()
	p.Wait(nil)
}
