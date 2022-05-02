package generator

type CompareChar struct {
	ch string
}

func (cc CompareChar) getKind() TestKind {
	return Char
}

func (cc CompareChar) accept(g GeneratorInterface) {
	g.VisitCompare(cc)
}
