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

	defaultNameDictionary = []string{
		"exusiai", "siege", "irit", "eyjafjalla", "angelina", "shining", "nightingale", "hoshiguma", "saria", "silverash",
		"skadi", "ch'en", "schwarz", "hellagur", "magallan", "mostima", "blaze", "aak", "nian", "ceobe", "bagpipe",
		"phantom", "weedy", "w", "роса", "suzuran", "thorns", "eunectes", "surtr", "blemishine", "mudrock", "rosmontis",
		"mountain", "archetto", "saga", "dusk", "ash", "passenger", "gladiia", "kal'tsit", "skadi-the-corrupting-heart",
		"carnelian", "pallas", "mizuki", "ch'en-the-holungday", "saileach", "fartooth", "flametail",
		"nearl-the-radiant-knight", "gnosis", "lee", "ling", "goldenglow", "fiammetta", "horn", "lumen", "irene",
		"specter-the-unchained", "ebenholz", "dorothy", "позёмка", "gavial-the-invincible", "młynar", "stainless",
		"vigil", "penance", "texas-the-omertosa", "reed-the-flame-shadow", "lin", "chongyue", "qiubai", "kirin-r-yato",
		"ines", "silence-the-paradigmatic", "ho'olheyak", "muelsyse", "executor-the-ex-foedere", "typhon",
		"swire-the-elegant-wit", "eyjafjalla-the-hvít-aska", "jessica-the-liberated", "hoederer", "lessing", "viviana",
		"virtuosa", "degenbrecher", "ray", "zuo-le", "shu", "ela", "ascalon", "civilight-eterna", "logos", "wiš'adel",
		"ulpianus", "nymph", "narantuya", "pepe", "marcille", "vina-victoria", "crownslayer", "vulpisfoglia",
		"lappland-the-decadenza", "thorns-the-lodestar", "blaze-the-igniting-spark", "yu", "entelechia", "necras",
		"mon3tr", "sankta-miksaparato", "lemuen", "exusiai-the-new-covenant", "tragodia", "leizi-the-thunderbringer",
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

func GenerateNameByDefaultDictionary(pk uint64) string {
	return GenerateName(pk, defaultNameDictionary)
}

func GenerateName(pk uint64, dictionary []string) (result string) {
	index := int(pk)%len(dictionary) - 1
	if index < 0 {
		index += len(dictionary)
	}

	suffix := IntegerAliasRoman(int(pk)/len(dictionary) + 1)

	return dictionary[index] + "." + suffix
}
