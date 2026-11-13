package dic

// Объект-одиночка справочника заголовков.
var singletonHeader Headers

// IHeader Интерфейс элемента справочника заголовков.
type IHeader interface {
	// IsEqual Истина, если заголовки эквивалентны.
	IsEqual(h IHeader) bool

	// String Заголовок в виде строки.
	String() string
}

// Структура объекта справочника заголовков.
type tHeader struct {
	header string // Заголовок.
}
