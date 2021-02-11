package msgfmt

import (
	"strings"
	"unicode"
)

type reader struct {
	*strings.Reader
}

func newReader(s string) *reader {
	return &reader{
		strings.NewReader(s),
	}
}

type runeCondition func(rune) bool

func equal(expected rune) runeCondition {
	return func(r rune) bool {
		return r == expected
	}
}

func anyOf(runes ...rune) (condition runeCondition, result func() rune) {
	var trigger rune

	return func(r rune) bool {
			for i := range runes {
				if r == runes[i] {
					trigger = r
					return true
				}
			}

			return false
		}, func() rune {
			return trigger
		}
}

func whitespace() runeCondition {
	return unicode.IsSpace
}

func keyword() runeCondition {
	first := true

	return func(r rune) bool {
		if first {
			first = false
			return unicode.IsLetter(r) || r == '_'
		}

		return unicode.IsDigit(r) || unicode.IsLetter(r) || r == '_'
	}
}

func (r *reader) ReadUntil(condition runeCondition) (prefix string, err error) {
	for {
		next, _, err := r.ReadRune()
		if err != nil || condition(next) {
			return prefix, err
		}

		prefix += string(next)
	}
}

func (r *reader) ReadWhile(condition runeCondition) (prefix string, err error) {
	return r.ReadUntil(func(next rune) bool {
		return !condition(next)
	})
}
