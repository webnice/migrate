package dic

import (
	"reflect"
	"strings"
)

// ParseMethod Разбор строки в объект HTTP метода запроса.
func ParseMethod(s string) IMethod {
	var (
		mdo *tMethod
		tmd *tMethod
		rv  reflect.Value
		n   int
	)

	s = strings.ToLower(strings.TrimSpace(s))
	rv = reflect.ValueOf(singletonMethod)
	for n = 0; n < rv.NumField(); n++ {
		tmd = rv.Field(n).Interface().(*tMethod)
		if strings.EqualFold(tmd.name, s) {
			mdo = tmd
		}
	}
	if mdo == nil {
		return nil
	}

	return mdo
}

// NewMethod Создание объекта под интерфейс IMethod, функция предназначена для формирования настраиваемых справочников.
func NewMethod(name string) IMethod {
	var (
		mdo *tMethod
		bit uint64
		rv  reflect.Value
		n   int
		tmd *tMethod
	)

	rv = reflect.ValueOf(singletonMethod)
	for n = 0; n < rv.NumField(); n++ {
		tmd = rv.Field(n).Interface().(*tMethod)
		if strings.EqualFold(tmd.name, name) {
			mdo = tmd
		}
		if tmd.bits > bit {
			bit = tmd.bits
		}
	}
	if mdo != nil {
		return mdo
	}
	mdo = &tMethod{
		name: strings.ToUpper(strings.TrimSpace(name)),
		bits: bit << 1,
	}

	return mdo
}
