//go:build tools

package ttyprogress

// Tools used by the Makefile, pinned by go.mod:
//   - ttyframes checks the behaviour of the examples
//     (targets behaviour and check-behaviour),
//   - mdref generates the README from doc/src, including
//     code fragments of the examples (targets doc and check-doc).
import (
	_ "github.com/mandelsoft/mdref"
	_ "github.com/mandelsoft/ttyframes/cmd/ttyframes"
)
