package dic

// Объект-одиночка справочника HTTP методов запросов.
var singletonMethod Methods

// IMethod Интерфейс HTTP методов запросов.
type IMethod interface {
	// IsEqual Истина, если HTTP методы эквивалентны, сравнивается только метод.
	IsEqual(m IMethod) bool

	// IsEqualFull Истина, если HTTP методы эквивалентны, сравнивается и метод и битовая маска.
	IsEqualFull(m IMethod) bool

	// String HTTP метод запроса.
	String() string
}

// Структура объекта HTTP метода запроса.
type tMethod struct {
	name string // HTTP метод запроса.
	bits uint64 // Биты HTTP метода.
}
