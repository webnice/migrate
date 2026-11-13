package dic

const mimeCharset = "charset"

// Объект-одиночка справочника MIME типов.
var singletonMimes Mimes

// IMime Интерфейс MIME типа.
type IMime interface {
	// Charset Кодовая страница MIME типа.
	Charset() (ret string)

	// IsEqual Истина, если MIME типы эквивалентны, без учёта атрибутов.
	IsEqual(m IMime) bool

	// IsEqualFull Истина, если MIME типы эквивалентны, включая атрибуты.
	IsEqualFull(m IMime) bool

	// Main Основной тип MIME, без подтипа и без опций.
	Main() string

	// Sub Подтип MIME, без опций.
	Sub() string

	// Mime Тип MIME с подтипом, но без опций.
	Mime() string

	// Opt Опция MIME типа с указанным ключом.
	Opt(key string) string

	// String Полностью весь MIME тип с подтипом и опциями.
	String() string
}

// Структура объекта MIME типа.
type tMime struct {
	main string            // Основной тип MIME.
	subt string            // Подтип MIME.
	opts map[string]string // Атрибуты MIME типа.
}
