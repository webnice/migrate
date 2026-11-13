package dic

import (
	"fmt"
	"sort"
	"strings"
)

// Генератор кода.
//go:generate go run mime_generate.go

// Mime Справочник MIME типов.
func Mime() Mimes { return singletonMimes }

// Charset Кодовая страница MIME типа.
func (mto *tMime) Charset() string { return mto.Opt(mimeCharset) }

// Main Основной тип MIME, без подтипа и без опций.
func (mto *tMime) Main() string { return mto.main }

// Sub Подтип MIME, без опций.
func (mto *tMime) Sub() string { return mto.subt }

// Mime Тип MIME с подтипом, но без опций.
func (mto *tMime) Mime() string {
	switch mto.subt != "" {
	case true:
		return fmt.Sprintf("%s/%s", mto.main, mto.subt)
	default:
		return mto.main
	}
}

// Opt Опция MIME типа с указанным ключом.
func (mto *tMime) Opt(key string) (ret string) {
	if _, ok := mto.opts[key]; !ok {
		return
	}
	ret = mto.opts[key]

	return
}

// String Полностью весь MIME тип с подтипом и опциями.
func (mto *tMime) String() (ret string) {
	const tpl = "%s; %s=%s"
	var (
		key string
		oky []string
		n   int
	)

	ret, oky = mto.Mime(), make([]string, 0, len(mto.opts))
	for key = range mto.opts {
		oky = append(oky, key)
	}
	sort.SliceStable(oky, func(i, j int) bool { return strings.Compare(oky[i], oky[j]) == -1 })
	for n = range oky {
		ret = fmt.Sprintf(tpl, ret, oky[n], mto.opts[oky[n]])
	}

	return
}

// IsEqual Истина, если MIME типы эквивалентны, без учёта атрибутов.
func (mto *tMime) IsEqual(m IMime) bool {
	if m == nil {
		return false
	}
	return strings.EqualFold(mto.Mime(), m.Mime())
}

// IsEqualFull Истина, если MIME типы эквивалентны, включая атрибуты.
func (mto *tMime) IsEqualFull(m IMime) bool {
	if m == nil {
		return false
	}
	return strings.EqualFold(mto.String(), m.String())
}
