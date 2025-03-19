package specs

type BarConfig struct {
	// Fill is the default character representing completed progress
	Fill rune
	// Head is the default character that moves when progress is updated
	Head rune
	// Empty is the default character that represents the empty progress
	Empty rune
	// LeftEnd is the default character in the left most part of the progress indicator
	LeftEnd rune
	// RightEnd is the default character in the right most part of the progress indicator
	RightEnd rune
}

func (c BarConfig) SetBrackets(b Brackets) BarConfig {
	c.LeftEnd = b.LeftEnd
	c.RightEnd = b.RightEnd
	return c
}

func (c BarConfig) SetBracketType(i int) BarConfig {
	if b, ok := BracketTypes[i]; ok {
		return c.SetBrackets(b)
	}
	return c
}

type Brackets struct {
	// LeftEnd is the default character in the left most part of the progress indicator
	LeftEnd rune
	// RightEnd is the default character in the right most part of the progress indicator
	RightEnd rune
}

func (b Brackets) Swap() Brackets {
	b.LeftEnd, b.RightEnd = b.RightEnd, b.LeftEnd
	return b
}

var (
	BracketTypes = map[int]Brackets{
		0: {'[', ']'},
		1: {'〚', '〛'},
		2: {'【', '】'},
		3: {'〖', '〗'},

		10: {'{', '}'},
		11: {'⦃', '⦄'},

		20: {'(', ')'},
		21: {'⦅', '⦆'},
		22: {'｟', '｠'},

		30: {'<', '>'},
		31: {'❮', '❯'},
		32: {'〈', '〉'},
		33: {'《', '》'},

		40: {'〔', '〕'},
		41: {'〘', '〙'},
		42: {'﹝', '﹞'},
		43: {'⦗', '⦘'},

		50: {'|', '|'},
		51: {'▏', '▏'},
		52: {'▐', '▌'},
	}

	// BarTypes describes predefined Bar configurations identified
	// by an integer.
	BarTypes = map[int]BarConfig{
		0: {
			Fill:     '=',
			Head:     '>',
			Empty:    '-',
			LeftEnd:  '[',
			RightEnd: ']',
		},
		1: {
			Fill: '═',
			Head: '▷',
			// Empty:    '┄',
			Empty:    '·',
			LeftEnd:  '▕',
			RightEnd: '▏',
		},
		2: {
			Fill:     '▬',
			Head:     '▶',
			Empty:    '┄',
			LeftEnd:  '▐',
			RightEnd: '▌',
		},

		10: {
			Fill:     '▒',
			Head:     '░',
			Empty:    '░',
			LeftEnd:  '▕',
			RightEnd: '▏',
		},
		11: {
			Fill:     '▓',
			Head:     '▒',
			Empty:    '▒',
			LeftEnd:  '▕',
			RightEnd: '▏',
		},
		12: {
			Fill:     '█',
			Head:     '▒',
			Empty:    '▒',
			LeftEnd:  '▕',
			RightEnd: '▏',
		},
	}
)
