package dic

import (
	"mime"
	"net/http"
	"path"
)

// Генератор кода.
//go:generate go run file_generate.go

// Регистрация полного справочника типов контента, расширений имён файлов.
func init() { fileAddExtensionType() }

// File Интерфейс справочника MIME типов расширений имён файлов.
func File() IFile { return singletonFile }

// Добавление всех MIME типов контента для расширений имён файлов.
func fileAddExtensionType() {
	const prefix = `.`
	var (
		mt string
		n  int
	)

	for mt = range mimeTypeExtension {
		for n = range mimeTypeExtension[mt] {
			_ = mime.AddExtensionType(prefix+mimeTypeExtension[mt][n], mt)
		}
	}
}

// MimeByFilename Определение MIME типа контента файла по имени файла.
// Если MIME тип не найден, вернётся nil объект.
func (feo *tFile) MimeByFilename(filename string) IMime {
	return feo.MimeByExtension(path.Ext(filename))
}

// MimeByExtension Определение MIME типа контента файла по расширению имени файла.
// Расширение имени файла должно начинаться с точки, пример: ".txt".
// Если MIME тип не найден, вернётся nil объект.
func (feo *tFile) MimeByExtension(extension string) (ret IMime) {
	var tmp string

	if tmp = mime.TypeByExtension(extension); tmp == "" {
		return
	}
	ret = ParseMime(tmp)

	return
}

// MimeByContent Определение MIME типа контента по срезу байт.
// Не стоит передавать весь файл, достаточно передать первые 512 байт или около того.
func (feo *tFile) MimeByContent(content []byte) (ret IMime) {
	const maxBuf = 512
	var (
		buf []byte
		n   int
	)

	buf = make([]byte, maxBuf)
	if n = copy(buf, content); n == 0 {
		return
	}
	ret = ParseMime(http.DetectContentType(buf[:n]))

	return
}
