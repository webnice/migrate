package dic

// Объект-одиночка справочника ошибок.
var singletonErrors Errors

// IError Интерфейс ошибок.
type IError interface {

	// CodeU64 Код ошибки с типом uint64.
	CodeU64() *Set[uint64]

	// CodeU32 Код ошибки с типом uint32.
	CodeU32() *Set[uint32]

	// CodeU16 Код ошибки с типом uint16.
	CodeU16() *Set[uint16]

	// CodeU8 Код ошибки с типом uint8.
	CodeU8() *Set[uint8]

	// CodeU Код ошибки с типом uint.
	CodeU() *Set[uint]

	// CodeI64 Код ошибки с типом int64.
	CodeI64() *Set[int64]

	// CodeI32 Код ошибки с типом int32.
	CodeI32() *Set[int32]

	// CodeI16 Код ошибки с типом int16.
	CodeI16() *Set[int16]

	// CodeI8 Код ошибки с типом int8.
	CodeI8() *Set[int8]

	// CodeI Код ошибки с типом int.
	CodeI() *Set[int]

	// TNew Назначение нового шаблона и списка аргументов ошибки.
	TNew(t string, arg ...string) IError

	// TArg Получение списка аргументов ошибки. Шаблон ошибки можно получить из якоря.
	TArg() []string

	// Anchor Якорь ошибки, он же шаблон ошибки или ошибка целиком, если аргументов нет.
	Anchor() error

	// Error Реализация интерфейса error.
	Error() string

	// String Представление ошибки в качестве строки, используется якорь и аргументы ошибки.
	String() string

	// Bytes Представление ошибки в качестве среза байт, используется якорь и аргументы ошибки.
	Bytes() []byte

	// Is Сравнение ошибок, вернётся истина если ошибки эквиваленты.
	Is(e error) (ret bool)

	// IsEqual Сравнение ошибок, вернётся истина если ошибки эквиваленты.
	IsEqual(e IError) bool

	// IsEqualCode Сравнение ошибок только по коду ошибки.
	IsEqualCode(e IError) (ret bool)

	// Bind Создание копии объекта IError, прикрепление к копии якоря и всех свойств, наполнение копии объекта ошибки
	// значениями аргументов и возврат объекта ошибки в качестве интерфейса стандартной ошибки.
	// Копии ошибок в таком виде, с разными значениями аргументов эквивалентны, если сравнивать функцией IsEqual().
	Bind(arg ...any) (err error)
}

// Errors Предопределённые ошибки.
type Errors struct {
	// Unknown Неизвестная ошибка.
	Unknown IError
}

// Структура объекта ошибки.
type tError struct {
	code         uint64   // Код ошибки.
	template     string   // Шаблон ошибки, он же якорь ошибки для сравнения эквивалентности.
	templateArgs []string // Описание аргументов ошибки.
	anchor       error    // Константа ошибки с фиксированным адресом.
	args         []any    // Значения аргументов возникшие в процессе работы приложения.
}
