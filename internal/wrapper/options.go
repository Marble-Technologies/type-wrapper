package wrapper

type Option func(*generator)

// Type sets type name to genarator.
func Type(typeName string) Option {
	return func(g *generator) {
		g.typ = typeName
	}
}

// Output sets output file path to genarator.
func Output(output string) Option {
	return func(g *generator) {
		g.output = output
	}
}

// Receiver sets receiver name to genarator.
func Receiver(receiver string) Option {
	return func(g *generator) {
		g.receiver = receiver
	}
}

// Lock sets lock field name to genarator.
func Lock(lock string) Option {
	return func(g *generator) {
		g.lock = lock
	}
}

// Implement io.Reader
func Reader(reader bool) Option {
	return func(g *generator) {
		g.reader = reader
	}
}

// Wrapper type name
func Wrapper(typ string) Option {
	return func(g *generator) {
		g.wrapperType = typ
	}
}

// Wrapper type name
func Interface(interfaceName string) Option {
	return func(g *generator) {
		g.interfaceName = interfaceName
	}
}
