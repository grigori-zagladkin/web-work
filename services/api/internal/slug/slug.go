package slug

import "strings"

const (
	MAX_LEN = 100
	A       = 'a'
	Z       = 'z'
	ZERO    = '0'
	NINE    = '9'
)

var DIVIDER byte = '-'
var TRANSLIT = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya"}

func isEnglishSymbol(runeSymbol rune) bool {
	return runeSymbol >= A && runeSymbol <= Z || runeSymbol >= ZERO && runeSymbol <= NINE
}

func Make(s string) string {
	var stringBuilder strings.Builder
	lastDash := true

	for _, runeSymbol := range strings.ToLower(s) {
		if translitedSymbol, ok := TRANSLIT[runeSymbol]; ok {
			stringBuilder.WriteString(translitedSymbol)
			if translitedSymbol != "" {
				lastDash = false
			}
			continue
		}
		switch {
		case isEnglishSymbol(runeSymbol):
			stringBuilder.WriteRune(runeSymbol)
			lastDash = false
		default:
			if !lastDash {
				stringBuilder.WriteByte(DIVIDER)
				lastDash = true
			}
		}
	}

	resultString := stringBuilder.String()
	if len(resultString) > MAX_LEN {
		resultString = resultString[:MAX_LEN]
	}
	return strings.Trim(resultString, string(DIVIDER))
}
