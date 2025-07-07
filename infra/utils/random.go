package utils

import (
	"math/rand/v2"
	"strings"
)

var (
	authCodeDictionary = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L",
		"M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}

	tokenDictionary = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L",
		"M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h",
		"i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}
)

func GenerateAuthCode(length int) string {
	seeds := make([]int, length)
	for i := 0; i < length; i++ {
		seeds[i] = rand.IntN(len(authCodeDictionary))
	}

	result := strings.Builder{}
	for _, seed := range seeds {
		result.WriteString(authCodeDictionary[seed])
	}

	return result.String()
}

func GenerateToken(length int, prefix string) string {
	seeds := make([]int, length)
	for i := 0; i < length; i++ {
		seeds[i] = rand.IntN(len(tokenDictionary))
	}

	result := strings.Builder{}
	result.WriteString(prefix)
	for _, seed := range seeds {
		result.WriteString(tokenDictionary[seed])
	}

	return result.String()
}
