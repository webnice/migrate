package dic

import (
	"reflect"
	"strings"
)

// ParseHeader Разбор строки в объект заголовка.
func ParseHeader(s string) IHeader {
	var (
		hro *tHeader
		thr *tHeader
		rv  reflect.Value
		n   int
	)

	s = strings.TrimSpace(s)
	rv = reflect.ValueOf(singletonHeader)
	for n = 0; n < rv.NumField(); n++ {
		if thr = rv.Field(n).Interface().(*tHeader); strings.EqualFold(thr.String(), s) {
			hro = thr
		}
	}
	if hro == nil {
		return nil
	}

	return hro
}

// NewHeader Создание объекта под интерфейс IHeader, функция предназначена для формирования настраиваемых справочников.
func NewHeader(header string) IHeader {
	var (
		hro *tHeader
		thr *tHeader
		rv  reflect.Value
		n   int
	)

	header = strings.TrimSpace(header)
	rv = reflect.ValueOf(singletonHeader)
	for n = 0; n < rv.NumField(); n++ {
		thr = rv.Field(n).Interface().(*tHeader)
		if strings.EqualFold(thr.String(), header) {
			hro = thr
		}
	}
	if hro != nil {
		return hro
	}
	hro = &tHeader{
		header: header,
	}

	return hro
}
