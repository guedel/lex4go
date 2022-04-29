package generator

type CompareLike struct {
	expression string
}

func (c CompareLike) getKind() TestKind {
	return Like
}

func (c CompareLike) accept(g GeneratorInterface) {
	g.VisitCompare(c)
}
