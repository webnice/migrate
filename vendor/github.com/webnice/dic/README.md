# dic

[![GoDoc](https://godoc.org/github.com/webnice/dic?status.png)](http://godoc.org/github.com/webnice/dic)
[![Coverage Status](https://coveralls.io/repos/github/webnice/dic/badge.svg?branch=master%0Av1)](https://coveralls.io/github/webnice/dic?branch=master%0Av1)

#### Описание

Библиотека справочников.
Реализованы следующие справочники:

* Mime - Справочник MIME типов.

* Method - Справочник HTTP методов запросов.

* Header - Справочник заголовков.

* Status - Справочник статусов HTTP ответов.

* File - Справочник MIME типов расширений файлов.

* Error - Справочник ошибок.

Справочники реализованы таким образом, чтобы их можно было дополнять пользовательскими данными через встраивание.
А получившийся пользовательский справочник обладал бы всеми свойствами справочника данной библиотеки и принимался
бы теми же функциями, которые используют исходный справочник.

#### Подключение
```bash
go get github.com/webnice/dic
```

### Пример дополнения справочника пользовательскими данными

```go
package main

import "fmt"

// CustomHeader Структура нового справочника с дополнительными пользовательскими заголовками.
type CustomHeader struct {
	Headers // Встраивание всех уже существующих заголовков.

	// MyCustomHeader1 Заголовок my-custom/header-1.
	MyCustomHeader1 IHeader

	// HeaderSecondCustom Заголовок header-second/custom.
	HeaderSecondCustom IHeader
}

// Реализация дополнительных пользовательских заголовков справочника.
var custom = &CustomHeader{
	Headers:            Header(),                          // Подключение всех существующих заголовков.
	MyCustomHeader1:    NewHeader("My-Custom/Header-1"),   // Заголовок my-custom/header-1.
	HeaderSecondCustom: NewHeader("Header-Second/Custom"), // Заголовок header-second/custom.
}

// Использование нового справочника.
func main() {
	var mch = custom.MyCustomHeader1
	var result = custom.AcceptPatch.IsEqual(mch)

	fmt.Printf("Является ли заголовок %q эквивалентным заголовку %q? Ответ: %t\n", mch, custom.AcceptPatch, result)
	hcc := custom.HeaderSecondCustom
	result = hcc.IsEqual(mch)
	fmt.Printf("Является ли заголовок %q эквивалентным заголовку %q? Ответ: %t\n", hcc, mch, result)
}
```

Результат:
```text
Является ли заголовок "My-Custom/Header-1" эквивалентным заголовку "Accept-Patch"? Ответ: false
Является ли заголовок "Header-Second/Custom" эквивалентным заголовку "My-Custom/Header-1"? Ответ: false
```
