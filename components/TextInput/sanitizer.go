package TextInput

import (
	"unicode"
	"unicode/utf8"
)

type Sanitizer struct{}

func NewSanitizer(opts ...func(*Sanitizer)) *Sanitizer {
	model := &Sanitizer{}

	for _, opt := range opts {
		opt(model)
	}

	return model
}

func (s *Sanitizer) Sanitize(runes []rune) []rune {
	sanitized_runes := []rune{}

	for idx := 0; idx < len(runes); idx++ {
		rn := runes[idx]

		switch {
		case rn == utf8.RuneError:
			// skip that one, we don't want it
		case unicode.IsControl(rn):
			// skip that one as well.
			// windows terminal quirk. When pressing ctrl by itself, it will send us an NUL rune
		default:
			sanitized_runes = append(sanitized_runes, runes[idx])
		}
	}

	return sanitized_runes
}
