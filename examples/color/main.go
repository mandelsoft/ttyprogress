package main

import (
	"github.com/fatih/color"
)

func main() {
	b := color.New(color.FgCyan, color.Bold)
	y := color.New(color.FgYellow)
	b.Printf("test %s end\n", y.Sprintf("yellow"))
}
