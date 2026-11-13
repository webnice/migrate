package dic

import (
	"errors"
	"reflect"
	"strings"
)

// Unbind Разыменовывает error, представленный стандартным интерфейсом ошибки, в объект IError.
// Если в error не содержится объекта удовлетворяющего интерфейсу IError, возвращается nil.
func (errs *Errors) Unbind(err error) (ierr IError) {
	var (
		bind *tError
		ok   bool
	)

	if ok = errors.As(err, &bind); ok {
		ierr = bind
	}

	return
}

// ParseError Анализ ошибки завёрнутой в стандартный интерфейс ошибки и извлечение IError.
// Если ошибка IError не найдена, создание новой ошибки с интерфейсом IError на базе ошибки error.
func (errs *Errors) ParseError(e error, db ...any) (ret IError) {
	var (
		bind *tError
		rv   reflect.Value
		n, d int
		ok   bool
	)

	// Список поиска.
	db = append([]any{errs}, db...)
	// Неизвестная ошибка.
	if e == nil {
		return singletonErrors.Unknown
	}
	if ret = errs.Unbind(e); ret != nil {
		return
	}
	// Поиск ошибки в справочниках.
	for d = 0; d < len(db) && bind == nil; d++ {
		rv = reflect.ValueOf(db[d])
		if rv.Kind() == reflect.Pointer {
			rv = rv.Elem()
		}
		for n = 0; n < rv.NumField(); n++ {
			if rv.Field(n).Kind() == reflect.Struct {
				continue
			}
			if !rv.Field(n).CanInterface() {
				continue
			}
			if bind, ok = rv.Field(n).Interface().(*tError); !ok || bind == nil || bind.anchor == nil {
				continue
			}
			if errors.Is(bind.anchor, e) || strings.EqualFold(bind.anchor.Error(), e.Error()) {
				break
			}
			bind = nil
		}
	}
	if bind != nil {
		ret = bind
		return
	}
	// Попытка извлечения корневой вложенной ошибки.
	bind = &tError{anchor: errors.Unwrap(e)}
	if bind.anchor == nil {
		bind.anchor = e
	}
	// Создание новой ошибки на базе error.
	ret = &tError{
		template:     e.Error(),
		templateArgs: []string{},
		anchor:       bind.anchor,
		args:         []any{},
	}

	return
}

// NewError Создание объекта под интерфейс IError,
// функция предназначена для формирования настраиваемых справочников.
func (errs *Errors) NewError(template string, arg ...string) IError {
	var erro *tError

	template = strings.TrimSpace(template)
	erro = &tError{
		template:     template,
		templateArgs: make([]string, 0, len(arg)),
		anchor:       errors.New(template),
	}
	erro.templateArgs = append(erro.templateArgs, arg...)

	return erro
}
