package main

import (
	"math/rand/v2"
	"os"
	"time"

	"github.com/mandelsoft/ttycolors"
	"github.com/mandelsoft/ttyprogress"
	"github.com/mandelsoft/ttyprogress/specs"
)

func main() {
	p := ttyprogress.For(os.Stdout)

	fmts := []ttycolors.FormatProvider{
		ttycolors.New(ttycolors.FmtRGBColor(90, 90, 90)),
		ttycolors.New(ttycolors.FmtRGBColor(110, 110, 110)),
		ttycolors.New(ttycolors.FmtRGBColor(140, 140, 140)),
		ttycolors.New(ttycolors.FmtRGBColor(170, 170, 170)),
		ttycolors.New(ttycolors.FmtRGBColor(140, 140, 140)),
		ttycolors.New(ttycolors.FmtRGBColor(110, 110, 110)),
	}
	{
		bar, _ := ttyprogress.NewPhasesSpinner(specs.NewFormatPhases("does-some-work", fmts...)).
			SetSpeed(2).
			AppendElapsed().Add(p)
		bar.Start()
		go func() {
			time.Sleep(time.Second * time.Duration(10+rand.Int()%20))
			bar.Close()
		}()
	}
	{
		bar, _ := ttyprogress.NewPhasesSpinner(specs.NewRollingTextPhases("does-some-work", 8)).
			SetSpeed(2).
			AppendElapsed().Add(p)
		bar.Start()
		go func() {
			time.Sleep(time.Second * time.Duration(20+rand.Int()%20))
			bar.Close()
		}()
	}
	{
		bar, _ := ttyprogress.NewPhasesSpinner(specs.NewBouncingTextPhases("work", 14)).
			SetSpeed(2).
			AppendElapsed().Add(p)
		bar.Start()
		go func() {
			time.Sleep(time.Second * time.Duration(20+rand.Int()%20))
			bar.Close()
		}()
	}
	{
		bar, _ := ttyprogress.NewPhasesSpinner(specs.NewBouncingTextFmtPhases("does-some-work", 3, ttycolors.New(ttycolors.FmtUnderline, ttycolors.FmtBlue))).
			SetSpeed(2).
			AppendElapsed().Add(p)
		bar.Start()
		go func() {
			time.Sleep(time.Second * time.Duration(20+rand.Int()%20))
			bar.Close()
		}()
	}
	{
		bar, _ := ttyprogress.NewPhasesSpinner(specs.NewRollingTextFmtPhases("does-some-work", 3, ttycolors.FmtBold)).
			SetSpeed(2).
			AppendElapsed().Add(p)
		bar.Start()
		go func() {
			time.Sleep(time.Second * time.Duration(20+rand.Int()%20))
			bar.Close()
		}()
	}

	p.Close()
	p.Wait(nil)
}
