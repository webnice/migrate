package dic

// Объект-одиночка справочника MIME типов файлов.
var singletonFile IFile = new(tFile)

// IFile Интерфейс справочника MIME типов файлов.
type IFile interface {
	// MimeByFilename Определение MIME типа контента файла по имени файла.
	// Если MIME тип не найден, вернётся nil объект.
	MimeByFilename(filename string) IMime

	// MimeByExtension Определение MIME типа контента файла по расширению имени файла.
	// Расширение имени файла должно начинаться с точки, пример: ".txt".
	// Если MIME тип не найден, вернётся nil объект.
	MimeByExtension(extension string) (ret IMime)

	// MimeByContent Определение MIME типа контента по срезу байт.
	// Не стоит передавать весь файл, достаточно передать первые 512 байт или около того.
	MimeByContent(content []byte) (ret IMime)
}

// Структура объекта справочника MIME типов файлов.
type tFile struct {
}
