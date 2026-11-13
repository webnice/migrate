package dic

import "strings"

// Генератор кода.
//go:generate go run status_generate.go

// Status Справочник статусов HTTP ответов.
func Status() Statuses { return singletonStatus }

// Code Числовой код статуса HTTP ответа.
func (sso *tStatus) Code() int { return sso.code }

// IsEqual Истина, если статусы эквивалентны.
func (sso *tStatus) IsEqual(s IStatus) (ret bool) {
	if s == nil {
		return false
	}
	if !strings.EqualFold(sso.String(), s.String()) {
		return
	}
	ret = sso.Code() == s.Code()

	return
}

// IsEqualCode Истина, если коды статусов HTTP ответов эквивалентны.
func (sso *tStatus) IsEqualCode(c int) bool { return sso.code == c }

// String Интерфейс Stringify. Статус HTTP ответа в виде строки.
func (sso *tStatus) String() string { return sso.status }

// Bytes Статус HTTP ответа в виде среза байт.
func (sso *tStatus) Bytes() []byte { return []byte(sso.String()) }
