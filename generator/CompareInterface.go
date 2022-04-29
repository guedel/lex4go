package generator

type TestKind uint

const (
	Eos = iota
	Char
	Charset
	Like
	Any
	In
	Or
	And
)

type AcceptGeneratorInterface interface {
	accept(g GeneratorInterface)
}

type GetKindInterface interface {
	getKind() TestKind
}

type CompareInterface interface {
	AcceptGeneratorInterface
	GetKindInterface
}
