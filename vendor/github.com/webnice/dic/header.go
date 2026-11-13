package dic

import "strings"

// Генератор кода.
//go:generate go run header_generate.go

// Header Справочник заголовков.
func Header() Headers { return singletonHeader }

// IsEqual Истина, если заголовки эквивалентны.
func (hro *tHeader) IsEqual(h IHeader) bool {
	if h == nil {
		return false
	}
	return strings.EqualFold(hro.String(), h.String())
}

// String Интерфейс Stringify. Заголовок в виде строки.
func (hro *tHeader) String() string { return hro.header }
