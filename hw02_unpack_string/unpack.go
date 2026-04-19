package hw02unpackstring

import (
	"errors"

	"strings"

	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(message string) (string, error) {

	var builder strings.Builder
	builder.Grow(len(message))
	escChrFlg := false
	var lastRune rune = 0

	for i:=0; i<len(message); {
		r, sz := utf8.DecodeRuneInString(message[i:])
		i += sz

		switch {
			case escChrFlg:
				escChrFlg=false
				lastRune=r
			case r == '\\':
				escChrFlg=true
				if lastRune !=0 {
					builder.WriteRune(lastRune)
					lastRune=0
				}
			case r>='0' && r<='9':
				var repNum int = int (r-'0');
				if lastRune == 0 {
					return "", ErrInvalidString
				}
				if repNum > 0 {
					builder.WriteString(strings.Repeat(string(lastRune),repNum))
				} 
				lastRune=0
			default:
				if lastRune !=0 {
					builder.WriteRune(lastRune)
				}
				lastRune=r
		}
	}

	if lastRune !=0 {
		builder.WriteRune(lastRune)
	}



	return builder.String(), nil
}
