package main

import (
	"fmt"
	"strings"
)

type AlgorithmGenerator struct {
	writer    CodeWriter
	testLevel int
}

func (g *AlgorithmGenerator) Quote(c string) string {
	return "\"" + g.Escape(c) + "\""
}

func (g *AlgorithmGenerator) Escape(c string) string {
	return strings.ReplaceAll(c, "\"", "\\\"")

}

func (g *AlgorithmGenerator) DoStartDocument(vars any) {
	if g.writer == (CodeWriter{}) {
		g.writer = newCodeWriter()
		g.writer.spacer = "\t"
	}
	const tpl string = `
// {{.Name}}	
// ----------------------------
// Création de {{.Author}}
// Créé le {{.CreationDate}}
// Modifié le {{.UpdateDate}}

Classe Scanner

	Structure DonnéesScanner
		var état: Chaine
		var ch: Caractère
		var token: Chaine
		var cnt: Entier
	Fin Structure
`
	g.writer.printTemplate(tpl, vars)
	g.writer.indentation = 2
}

func (g *AlgorithmGenerator) DoEndDocument(vars any) {
	g.writer.unindent()
	g.writer.println("Fin Classe")
}

func (g *AlgorithmGenerator) DoGenerateProlog(vars any) {
	const tpl string = `
	Procédure Scanner(s: chaine)
		var ancien_état: Chaine = ""
		var p: Entier = 0
		var d: DonnéesScanner

		Avec d
			.état = "0"
			.ch = ''
			.token = ""
			.cnt = 0
		Fin Avec

		s <- s + chr(0)

		Pour p de 1 à longueur_chaine(s) Faire
			Si d.état = ancien_état Alors
				d.cnt <- d.cnt + 1
			Sinon
				d.cnt <- 0
			Fin Si

			ancien_état <- d.état
			d.ch <- s[p]

			Selon d.état
`
	g.writer.printTemplate(tpl, vars)
	g.writer.indentation = 5
}

func (g *AlgorithmGenerator) DoGenerateEpilog(finalState string) {
	const tpl string = `
				Sinon
					Lève Exception_UnexpectedChar
				Fin Si
			Fin Selon
		Fin Pour
		Si état <> "{{.finalState}}" Alors
			Lève Exception_UnexpectedEnd
		Fin Si
	Fin Procédure
`
	var vars = make(map[string]string)
	vars["finalState"] = g.Escape(finalState)
	g.writer.printTemplate(tpl, vars)
	g.writer.nl()
	g.writer.indentation = 1
}

func (g *AlgorithmGenerator) DoClosePreviousIf() {
	w := &(g.writer)
	w.unindent()
	w.println("Sinon").indent()
	w.println("Lève Exception_UnexcpectedChar").unindent()
	w.println("Fin Si").nl()
}

func (g *AlgorithmGenerator) DoStartNewState(state string) {
	w := &(g.writer)
	w.unindent()
	w.println("Cas " + g.Quote(state) + ":")
	w.indent()
	w.print("Si ")
}

func (g *AlgorithmGenerator) DoElseIf() {
	w := &(g.writer)
	w.unindent()
	w.print("Sinon si ")
}

func (g *AlgorithmGenerator) DoTestLike(expr string) {
	g.writer.println("d.ch Comme " + g.Quote(expr) + " Alors").indent()
}

func (g *AlgorithmGenerator) DoTestCharset(charset string) {
	var test string
	switch strings.ToUpper(charset) {
	case "ALPHA":
		test = "estAlpha(ch)"
	case "BLANK":
		test = "estBlanc(ch)"
	case "CONTROL":
		test = "estControle(ch)"
	case "DIGIT":
		test = "estChiffre(ch)"
	case "GRAPH":
		test = "estGraphe(ch)"
	case "LOWER":
		test = "estMinuscule(ch)"
	case "PRINT":
		test = "estImprimable(ch)"
	case "PUNCT":
		test = "estPonctuation(ch)"
	case "UPPER":
		test = "estMajuscule(ch)"
	case "XDIGIT":
		test = "estChiffreHexa(ch)"
	}
	g.writer.print(test)
}

func (g *AlgorithmGenerator) DoTestAny() {
	g.writer.println("vrai Alors").indent()
}

func (g *AlgorithmGenerator) DoMaxLoops(maxloops int) {
	w := &(g.writer)
	w.println(fmt.Sprintf("Si d.cnt > %d Alors", maxloops)).indent()
	w.println("Lève Exception_UnexpectedChar").unindent()
	w.println("Fin Si")
}

func (g *AlgorithmGenerator) DoAction(action string, transition string, useLoops bool, testMode bool) {
	if testMode {
		g.writer.println(fmt.Sprintf("Appel de '%s' pour la transition '%s'", g.Escape(action), g.Escape(transition)))
	} else {
		g.writer.println(fmt.Sprintf("OnAction %s, %s, d", g.Quote(action), g.Quote(transition)))
	}
}

func (g *AlgorithmGenerator) DoResetToken() {
	g.writer.println("d.token <- \"\"")
}

func (g *AlgorithmGenerator) DoAddToToken() {
	g.writer.println("d.token <- concatène(d.token, d.ch)")
}

func (g *AlgorithmGenerator) DoNewState(state string) {
	g.writer.println("d.état <- " + g.Quote(state))
}

func (g *AlgorithmGenerator) DoPrototype(actions []string) {
	w := &(g.writer)
	w.println("Procédure OnAction(sTransition, sAction: Chaine; data: DonnéesScanne)").indent()
	w.println("Selon sAction").indent()
	for _, action := range actions {
		w.println("Cas " + g.Quote(action) + ":").indent()
		w.println("// votre code ici").nl().unindent()
	}
	w.unindent()
	w.println("Fin Selon").unindent()
	w.println("Fin Procédure")
}

func (g *AlgorithmGenerator) VisitCompare(c CompareInterface) {
	w := &(g.writer)
	switch c := c.(type) {
	case *CompareEos:
		w.print("Code(d.ch) = 0")
	case CompareAny:
		w.print("vrai")
	case *CompareChar:
		w.print("d.ch = '" + c.ch + "'")
	case *CompareCharset:
		g.DoTestCharset(c.Name)
	case *CompareLike:
		w.print("d.ch Commme " + g.Quote(c.Expression))
	case *CompareAnd:
		w.print("(")
		g.testLevel++
		for index, child := range c.childs {
			if index > 0 {
				w.print(" Et ")
			}
			child.accept(g)
		}
		g.testLevel--
		w.print(")")
	case *CompareOr:
		w.print("(")
		g.testLevel++
		for index, child := range c.childs {
			if index > 0 {
				w.print(" Ou ")
			}
			child.accept(g)
		}
		g.testLevel--
		w.print(")")
		/*
			case In:
				w.print("Dans (")
				w.print(")")
		*/
	}

	if g.testLevel == 0 {
		w.println(" Alors").indent()
	}
}
