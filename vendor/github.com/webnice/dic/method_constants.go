package dic

// Methods Структура справочника HTTP методов запросов.
type Methods struct {
	// Get HTTP метод GET.
	Get IMethod
	// Post HTTP метод POST.
	Post IMethod
	// Put HTTP метод PUT.
	Put IMethod
	// Delete HTTP метод DELETE.
	Delete IMethod
	// Connect HTTP метод CONNECT.
	Connect IMethod
	// Head HTTP метод HEAD.
	Head IMethod
	// Patch HTTP метод PATCH.
	Patch IMethod
	// Options HTTP метод OPTIONS.
	Options IMethod
	// Trace HTTP метод TRACE.
	Trace IMethod
	// Stub HTTP метод STUB.
	Stub IMethod
}

func init() {
	singletonMethod = Methods{
		Get:     &tMethod{name: "GET", bits: 1 << 0},
		Post:    &tMethod{name: "POST", bits: 1 << 1},
		Put:     &tMethod{name: "PUT", bits: 1 << 2},
		Delete:  &tMethod{name: "DELETE", bits: 1 << 3},
		Connect: &tMethod{name: "CONNECT", bits: 1 << 4},
		Head:    &tMethod{name: "HEAD", bits: 1 << 5},
		Patch:   &tMethod{name: "PATCH", bits: 1 << 6},
		Options: &tMethod{name: "OPTIONS", bits: 1 << 7},
		Trace:   &tMethod{name: "TRACE", bits: 1 << 8},
		Stub:    &tMethod{name: "STUB", bits: 1 << 9},
	}
}
