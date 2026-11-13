package dic

import (
	"reflect"
	"strings"
)

// ParseStatusString Разбор строки в объект статуса HTTP ответа.
func ParseStatusString(s string) IStatus {
	var (
		sso *tStatus
		tss *tStatus
		rv  reflect.Value
		n   int
	)

	s = strings.TrimSpace(s)
	rv = reflect.ValueOf(singletonStatus)
	for n = 0; n < rv.NumField(); n++ {
		if tss = rv.Field(n).Interface().(*tStatus); strings.EqualFold(tss.String(), s) {
			sso = tss
		}
	}
	if sso == nil {
		return nil
	}

	return sso
}

// ParseStatusCode Разбор числового кода в объект статуса HTTP ответа.
func ParseStatusCode(code int) IStatus {
	var (
		sso *tStatus
		tss *tStatus
		rv  reflect.Value
		n   int
	)

	rv = reflect.ValueOf(singletonStatus)
	for n = 0; n < rv.NumField(); n++ {
		if tss = rv.Field(n).Interface().(*tStatus); tss.Code() == code {
			sso = tss
		}
	}
	if sso == nil {
		return nil
	}

	return sso
}

// NewStatus Создание объекта под интерфейс IStatus, функция предназначена для формирования настраиваемых справочников.
func NewStatus(status string, code int) IStatus {
	var (
		sso *tStatus
		tss *tStatus
		rv  reflect.Value
		n   int
	)

	status = strings.TrimSpace(status)
	rv = reflect.ValueOf(singletonStatus)
	for n = 0; n < rv.NumField(); n++ {
		tss = rv.Field(n).Interface().(*tStatus)
		switch tss.Code() == code {
		case true:
			sso = tss
		default:
			if sso == nil && strings.EqualFold(tss.String(), status) {
				sso = tss
			}
		}
	}
	if sso != nil {
		return sso
	}
	sso = &tStatus{
		status: status,
		code:   code,
	}

	return sso
}
