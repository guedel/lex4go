package main

type CompareAny struct {
}

func (ca CompareAny) getKind() TestKind {
	return Any
}

func (ca CompareAny) accept(g GeneratorInterface) {
	g.VisitCompare(ca)
}
