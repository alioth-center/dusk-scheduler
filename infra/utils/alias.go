package utils

import "strings"

var symbols = []struct {
	value  int
	symbol string
}{
	{1000, "m"}, {900, "cm"}, {500, "d"}, {400, "cd"},
	{100, "c"}, {90, "xc"}, {50, "l"}, {40, "xl"},
	{10, "x"}, {9, "ix"}, {5, "v"}, {4, "iv"},
	{1, "i"},
}

func IntegerAliasRoman(number int) (roman string) {
	result := strings.Builder{}
	for _, symbol := range symbols {
		for number >= symbol.value {
			result.WriteString(symbol.symbol)
			number -= symbol.value
		}

		if number == 0 {
			break
		}
	}

	return result.String()
}
