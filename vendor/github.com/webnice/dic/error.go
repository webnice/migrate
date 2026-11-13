package dic

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

func init() {
	singletonErrors = Errors{
		Unknown: NewError("неизвестная ошибка").CodeI64().Set(math.MinInt64),
	}
}

// Error Справочник ошибок.
func Error() Errors { return singletonErrors }

// CodeU64 Код ошибки с типом uint64.
func (eso *tError) CodeU64() *Set[uint64] { return newSet[uint64](eso) }

// CodeU32 Код ошибки с типом uint32.
func (eso *tError) CodeU32() *Set[uint32] { return newSet[uint32](eso) }

// CodeU16 Код ошибки с типом uint16.
func (eso *tError) CodeU16() *Set[uint16] { return newSet[uint16](eso) }

// CodeU8 Код ошибки с типом uint8.
func (eso *tError) CodeU8() *Set[uint8] { return newSet[uint8](eso) }

// CodeU Код ошибки с типом uint.
func (eso *tError) CodeU() *Set[uint] { return newSet[uint](eso) }

// CodeI64 Код ошибки с типом int64.
func (eso *tError) CodeI64() *Set[int64] { return newSet[int64](eso) }

// CodeI32 Код ошибки с типом int32.
func (eso *tError) CodeI32() *Set[int32] { return newSet[int32](eso) }

// CodeI16 Код ошибки с типом int16.
func (eso *tError) CodeI16() *Set[int16] { return newSet[int16](eso) }

// CodeI8 Код ошибки с типом int8.
func (eso *tError) CodeI8() *Set[int8] { return newSet[int8](eso) }

// CodeI Код ошибки с типом int.
func (eso *tError) CodeI() *Set[int] { return newSet[int](eso) }

// TNew Назначение нового шаблона и списка аргументов ошибки.
func (eso *tError) TNew(t string, arg ...string) IError {
	if t != "" {
		eso.anchor = errors.New(t)
	}
	eso.templateArgs = make([]string, 0, len(arg))
	eso.templateArgs = append(eso.templateArgs, arg...)
	return eso
}

// TArg Получение списка аргументов шаблона ошибки.
func (eso *tError) TArg() []string { return eso.templateArgs }

// Anchor Якорь ошибки, он же шаблон ошибки или ошибка целиком, если аргументов нет.
func (eso *tError) Anchor() error { return eso.anchor }

// Error Реализация интерфейса error.
func (eso *tError) Error() string { return eso.String() }

// String Представление ошибки в качестве строки, используется якорь и аргументы ошибки.
func (eso *tError) String() string { return fmt.Sprintf(eso.template, eso.args...) }

// Bytes Представление ошибки в качестве среза байт, используется якорь и аргументы ошибки.
func (eso *tError) Bytes() []byte { return []byte(eso.String()) }

// Is Сравнение ошибок, вернётся истина если ошибки эквиваленты.
func (eso *tError) Is(e error) (ret bool) {
	var ierr IError

	if ierr = singletonErrors.ParseError(e, eso); ierr != nil {
		ret = eso.IsEqual(ierr)
	}

	return
}

// IsEqual Сравнение ошибок, вернётся истина если ошибки эквиваленты.
func (eso *tError) IsEqual(e IError) (ret bool) {
	if e == nil {
		return
	}
	// Сравнение стандартной библиотекой.
	ret = errors.Is(eso.Anchor(), e.Anchor())
	// Сравнение по якорю через эквивалент строк.
	if !ret {
		ret = strings.EqualFold(eso.Anchor().Error(), e.Anchor().Error())
	}

	return
}

// IsEqualCode Сравнение ошибок только по коду ошибки.
func (eso *tError) IsEqualCode(e IError) (ret bool) {
	ret = eso.CodeU64().Get() == e.CodeU64().Get()
	return
}

// Bind Создание копии объекта IError, прикрепление к копии якоря и всех свойств, наполнение копии объекта ошибки
// значениями аргументов и возврат объекта ошибки в качестве интерфейса стандартной ошибки.
// Копии ошибок в таком виде, с разными значениями аргументов эквивалентны, если сравнивать функцией IsEqual().
func (eso *tError) Bind(arg ...any) (err error) {
	var (
		bind *tError
		n    int
	)

	bind = &tError{
		code:         eso.code,
		template:     eso.template,
		templateArgs: eso.templateArgs,
		anchor:       eso.anchor,
		args:         make([]any, len(eso.templateArgs)),
	}
	for n = range bind.args {
		if len(arg) > n {
			bind.args[n] = arg[n]
			continue
		}
		if len(eso.templateArgs) > n {
			bind.args[n] = eso.templateArgs[n]
			continue
		}
	}

	return bind
}

// NewError Создание объекта под интерфейс IError,
// функция предназначена для формирования настраиваемых справочников.
func NewError(template string, arg ...string) IError {
	return singletonErrors.NewError(template, arg...)
}

// ParseError Анализ ошибки завёрнутой в стандартный интерфейс ошибки и извлечение IError.
// Если ошибка IError не найдена, создание новой ошибки с интерфейсом IError на базе ошибки error.
func ParseError(e error) IError { return singletonErrors.ParseError(e) }
