package dic

// Модель дженерика типов значений.
type modelSet interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

// Set Структура объекта доступа к коду.
type Set[T modelSet] struct {
	parent *tError
}

// Конструктор объекта.
func newSet[T modelSet](p *tError) *Set[T] {
	return &Set[T]{parent: p}
}

// Set Присвоение значения.
func (sto *Set[T]) Set(number T) IError {
	sto.parent.code = uint64(number)
	return sto.parent
}

// Get Возвращение значения.
func (sto *Set[T]) Get() (ret T) {
	ret = T(sto.parent.code)
	return
}
