package dic

import "strings"

// ParseMime Разбор строки в объект MIME типа.
func ParseMime(s string) IMime {
	const keySlash, keyEqual, keySemicolon = "/", "=", ";"
	var (
		clean func(string) string
		mto   *tMime
		sss   []string
		tmp   []string
		n     int
	)

	clean = func(s string) string { return strings.ToLower(strings.TrimSpace(s)) }
	mto, sss = &tMime{opts: make(map[string]string)}, strings.Split(s, keySemicolon)
	switch tmp = strings.Split(sss[0], keySlash); len(tmp) {
	case 1:
		mto.main = clean(tmp[0])
	default:
		mto.main = clean(tmp[0])
		mto.subt = clean(tmp[1])
	}
	if len(sss) < 2 {
		return mto
	}
	for n = 1; n < len(sss); n++ {
		switch tmp = strings.Split(sss[n], keyEqual); len(tmp) {
		case 2:
			mto.opts[clean(tmp[0])] = clean(tmp[1])
		}
	}

	return mto
}

// NewMime Создание объекта под интерфейс IMime, функция предназначена для формирования настраиваемых справочников.
func NewMime(main, sub string, opts map[string]string) IMime {
	var (
		mto *tMime
		key string
	)

	mto = &tMime{
		main: main,
		subt: sub,
		opts: make(map[string]string),
	}
	for key = range opts {
		if key == "" || opts[key] == "" {
			continue
		}
		mto.opts[key] = opts[key]
	}

	return mto
}
