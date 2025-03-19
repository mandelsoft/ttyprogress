package specs

import (
	"maps"

	"github.com/briandowns/spinner"
)

// SpinnerTypes predefined spinner types.
// Most of them are taken from [spinner.CharSets] (github.com/briandowns/spinner).
var SpinnerTypes = maps.Clone(spinner.CharSets)

func init() {
	SpinnerTypes[1000] = []string{"███▒▒▒▒▒▒▒", "▒███▒▒▒▒▒▒", "▒▒███▒▒▒▒▒", "▒▒▒███▒▒▒▒", "▒▒▒▒███▒▒▒", "▒▒▒▒▒███▒▒", "▒▒▒▒▒▒███▒", "▒▒▒▒▒▒▒███", "█▒▒▒▒▒▒▒██", "██▒▒▒▒▒▒▒█"}
	SpinnerTypes[1001] = []string{"███▒▒▒▒▒▒▒", "▒███▒▒▒▒▒▒", "▒▒███▒▒▒▒▒", "▒▒▒███▒▒▒▒", "▒▒▒▒███▒▒▒", "▒▒▒▒▒███▒▒", "▒▒▒▒▒▒███▒", "▒▒▒▒▒▒▒███", "▒▒▒▒▒▒███▒", "▒▒▒▒▒███▒▒", "▒▒▒▒███▒▒▒", "▒▒███▒▒▒▒▒", "▒███▒▒▒▒▒▒"}
	SpinnerTypes[1002] = []string{"⋮", "⋰", "⋯", "⋱"}
	SpinnerTypes[1003] = []string{"✶", "✷", "✸", "✷"}
	SpinnerTypes[1004] = []string{"𝄖", "𝄗", "𝄘", "𝄙", "𝄛", "𝄙", "𝄘", "𝄗", "𝄖"}
	SpinnerTypes[1005] = []string{"▢", "▣"}
	SpinnerTypes[1006] = []string{"◇", "◈"}
	SpinnerTypes[1007] = []string{"◇", "◈", "◆"}
}
