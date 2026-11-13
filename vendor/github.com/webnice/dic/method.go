package dic

import "strings"

// Method Справочник HTTP методов запросов.
func Method() Methods { return singletonMethod }

// IsEqual Истина, если HTTP методы эквивалентны.
func (mdo *tMethod) IsEqual(m IMethod) bool {
	if m == nil {
		return false
	}
	return strings.EqualFold(mdo.name, m.(*tMethod).name)
}

// IsEqualFull Истина, если HTTP методы эквивалентны, сравнивается и метод и битовая маска.
func (mdo *tMethod) IsEqualFull(m IMethod) (ret bool) {
	if m == nil {
		return false
	}
	ret = mdo.IsEqual(m) && mdo.bits == m.(*tMethod).bits

	return
}

// String Интерфейс Stringify.
func (mdo *tMethod) String() string { return mdo.name }
