package main

import (
	"errors"
	"testing"

	"github.com/webnice/dic"
)

func TestErrorsSingleton(t *testing.T) {
	const errSingletonBroken = "функция Errors(), функция повреждена"

	if Errors() != errSingleton {
		t.Fatal(errSingletonBroken)
	}
}

// TestErrors Тестирование всех ошибок пакета - корректность текста ошибки и эквивалентность справочнику.
//
//goland:noinspection DuplicatedCode
func TestErrors(t *testing.T) {
	const (
		errBoom = "тестовая ошибка"
	)
	var (
		base  error
		items = []struct {
			name string
			got  error
			ref  dic.IError
			text string
		}{
			{"UnknownDriver", Errors().UnknownDriver.Bind("oracle"), Errors().UnknownDriver, `"oracle": драйвер не поддерживается.`},
		}
	)

	for n := range items {
		if items[n].got.Error() != items[n].text {
			t.Errorf("функция ошибки %q, текст ошибки: %q, ожидался: %q", items[n].name, items[n].got.Error(), items[n].text)
		}
		if !items[n].ref.Is(items[n].got) {
			t.Errorf("функция ошибки %q, ошибка не эквивалентна справочнику", items[n].name)
		}
	}
	// Ошибки с обёрнутой исходной ошибкой: корректность текста и эквивалентность справочнику.
	base = errors.New(errBoom)
	for _, wrapped := range []struct {
		name string
		got  error
		ref  dic.IError
	}{
		{"ConnectDatabase", Errors().ConnectDatabase.Bind(base), Errors().ConnectDatabase},
		{"RunCommand", Errors().RunCommand.Bind("up", base), Errors().RunCommand},
	} {
		if !wrapped.ref.Is(wrapped.got) {
			t.Errorf("функция ошибки %q, ошибка не эквивалентна справочнику", wrapped.name)
		}
	}
}
