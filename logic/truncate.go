package logic

import (
	"unicode"
	"unicode/utf8"
)

func TruncateByWords(s string, maxWords int) string {
	processdWords := 0
	wordStarted := false
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRuneInString(s[i:])
		if !isSeparator(r) {
			i += width
			wordStarted = true
			continue
		}

		if !wordStarted {
			i += width
			continue
		}

		wordStarted = false
		processdWords++
		if processdWords == maxWords {
			const ending = "..."
			if (i + len(ending)) >= len(s) {
				return s
			}
			return s[:i] + ending
		}
		i += width
	}
	return s
}

// 判断字符是不是分隔符
func isSeparator(r rune) bool {
	// ASCII字符和下划线不是分隔符
	if r <= 0x7F {
		switch {
			case '0' <= r && r <= '9':
				return false
			case 'a' <= r && r <= 'z':
				return false
			case 'A' <= r && r <= 'Z':
				return false
			case r == '_':
				return false
			default:
				return true
		}
	}
	//字母和数组不是分隔符
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	// 空格是分隔符
	return unicode.IsSpace(r)
}