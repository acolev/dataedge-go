package dataedge

import (
	"strings"
)

var translitMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v",
	'г': "g", 'д': "d", 'е': "e",
	'ё': "e", 'ж': "zh", 'з': "z",
	'и': "i", 'й': "y", 'к': "k",
	'л': "l", 'м': "m", 'н': "n",
	'о': "o", 'п': "p", 'р': "r",
	'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "c",
	'ч': "ch", 'ш': "sh", 'щ': "sch",
	'ь': "'", 'ы': "y", 'ъ': "'",
	'э': "e", 'ю': "yu", 'я': "ya",

	'А': "A", 'Б': "B", 'В': "V",
	'Г': "G", 'Д': "D", 'Е': "E",
	'Ё': "E", 'Ж': "Zh", 'З': "Z",
	'И': "I", 'Й': "Y", 'К': "K",
	'Л': "L", 'М': "M", 'Н': "N",
	'О': "O", 'П': "P", 'Р': "R",
	'С': "S", 'Т': "T", 'У': "U",
	'Ф': "F", 'Х': "H", 'Ц': "C",
	'Ч': "Ch", 'Ш': "Sh", 'Щ': "Sch",
	'Ь': "'", 'Ы': "Y", 'Ъ': "'",
	'Э': "E", 'Ю': "Yu", 'Я': "Ya",
}

// Transliterate converts Cyrillic characters to Latin.
func Transliterate(text string) string {
	var builder strings.Builder
	for _, r := range text {
		if val, ok := translitMap[r]; ok {
			builder.WriteString(val)
		} else {
			builder.WriteRune(r)
		}
	}
	result := builder.String()

	// Check casing of first character to match PHP logic
	// PHP: $case = preg_match('~^\p{Lu}~u', $str) ? 'upper' : 'lower';
	// PHP: if ($case == "upper") return ucfirst($str);
	// The map already handles case for the first character if it was in the map.
	// But if the first char wasn't in the map (e.g. Latin), it stays as is.
	// PHP's logic seems to force ucfirst if the *transliterated* string starts with uppercase?
	// Or rather, checking if the *original* string started with uppercase?
	// `preg_match('~^\p{Lu}~u', $str)` checks the *transliterated* string ($str).
	// If the first char of result is uppercase, it runs `ucfirst`... which does nothing if it's already uppercase?
	// Wait, PHP `ucfirst` capitalizes the first char.
	// If it's already upper, it stays upper.
	// So that logic in PHP seems redundant unless `strtr` produces lowercase for uppercase input?
	// The PHP map has 'А' => 'A'. So it produces uppercase.
	// So `ucfirst` is redundant.
	// unless...
	// Ah, there is a commented out line in PHP: `$str = trim(preg_replace... strtolower($str))...`
	// Maybe that was old logic.
	// The current PHP code essentially just does replacement.

	// However, let's look closely at `preg_match`. It checks `$str` (the result).
	// So if result starts with Upper, it returns `ucfirst($str)`.
	// Since it's already upper, it's a no-op.
	// If it starts with lower, it returns `$str`.
	// So... this logic effectively does nothing extra besides the map replacement?
	// Verified. I will just stick to the map replacement.

	return result
}
