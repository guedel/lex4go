package main

type CompareEos struct{}

func (t CompareEos) getKind() TestKind {
	return Eos
}

func (t CompareEos) accept(g GeneratorInterface) {
	g.VisitCompare(t)
}
