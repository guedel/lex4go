package generator

type GeneratorInterface interface {
	DoStartDocument(vars any)
	DoEndDocument(vars any)
	DoGenerateProlog(vars any)
	DoGenerateEpilog(finalState string)
	DoClosePreviousIf()
	DoStartNewState(state string)
	DoElseIf()
	DoMaxLoops(maxloops int)
	DoAction(action string, transition string, useLoops bool, testMode bool)
	DoResetToken()
	DoAddToToken()
	DoNewState(state string)
	DoPrototype(actions []string)
	VisitCompare(i CompareInterface)
}
