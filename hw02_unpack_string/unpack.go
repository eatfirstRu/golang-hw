package hw02_unpack_string

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inStr string) (string, error) {
	isEscSymb := false
	outStr := ""

	if inStr == "" {
		return "", nil
	}

	inRune := []rune(inStr)

	for i := 0; i < len(inRune); i++ {

		if inRune[i] == '\\' && !isEscSymb {
			isEscSymb = true
			continue
		}

		isCurDigit := unicode.IsDigit(inRune[i])
		isCurLetter := unicode.IsLetter(inRune[i])
		isCurSpace := unicode.IsSpace(inRune[i])

		// обработчик последнего символа
		if i == len(inRune)-1 {
			if isCurLetter || isCurSpace || isEscSymb && isCurDigit {
				outStrAdd(&outStr, inRune[i], 1)
				return outStr, nil
			} else {
				return "", ErrInvalidString
			}
		}

		isNextLetter := unicode.IsLetter(inRune[i+1])
		isNextSpace := unicode.IsSpace(inRune[i+1])
		isNextDigit := unicode.IsDigit(inRune[i+1])
		isNextBackSlash := inRune[i+1] == '\\'

		// если текущий символ буква или спец символ или перед ним стоит isEscSymb=\\
		if isCurLetter || isCurSpace || isEscSymb {
			//смотрим следующий символ
			//цифра, печатаем текущий символ cntRepeat раз
			if isNextDigit {
				cntRepeat, _ := strconv.Atoi(string(inRune[i+1]))
				outStrAdd(&outStr, inRune[i], cntRepeat)
				i++
				isEscSymb = false
				continue
			}
			//буква или спец символ, печатаем текущий символ 1 раз
			if isNextLetter || isNextSpace || isNextBackSlash {
				outStrAdd(&outStr, inRune[i], 1)
				isEscSymb = false
			}
		} else {
			return "", ErrInvalidString
		}
	}
	return outStr, nil
}

func outStrAdd(outStr *string, curRune rune, cntRepeat int) {
	for i := 1; i <= cntRepeat; i++ {
		*outStr += string(curRune)
	}
}
