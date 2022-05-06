package generator

import (
	"fmt"
	"strings"
)

type GoGenerator struct {
	writer    CodeWriter
	testLevel int
	functions []string
}

func (g *GoGenerator) Quote(c string) string {
	return "\"" + g.Escape(c) + "\""
}

func (g *GoGenerator) Escape(c string) string {
	return strings.ReplaceAll(c, "\"", "\\\"")

}

func (g *GoGenerator) DoStartDocument(vars any) {
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
// ----------------------------
// {{.Description}}

package Scanner

import (
	"fmt"
	"strings"
)	

type ScannerData struct {
	state string
	ch rune
	token string
	cnt int
}

`
	g.writer.printTemplate(tpl, vars)
	g.writer.indentation = 1
}

func (g *GoGenerator) DoEndDocument(vars any) {
	// Ici on ajoute les fonctions demandées
	w := &(g.writer)
	for _, fn := range g.functions {
		switch fn {
		case "isAlpha":
			w.nl()
			w.println("func isAlpha(ch rune) bool {").indent()
			w.println("return isLower(ch) || isUpper(ch)").unindent()
			w.println("}")
			g.functions = AddUnique(g.functions, "isLower")
			g.functions = AddUnique(g.functions, "isUpper")
		case "isAlNum":
			w.nl()
			w.println("func isAlNum(ch rune) bool {").indent()
			w.println("return isAlpha(ch) || isDigit(ch)").unindent()
			w.println("}")
			g.functions = AddUnique(g.functions, "isAlpha")
			g.functions = AddUnique(g.functions, "isDigit")

		case "isDigit":
			w.nl()
			w.println("func isDigit(ch rune) bool {").indent()
			w.println("return ch >= 48 && ch <= 57").unindent()
			w.println("}")

		case "isUpper":
			w.nl()
			w.println("func isUpper(ch rune) bool {").indent()
			w.println("return ch >= 65 && ch <= 90").unindent()
			w.println("}")

		case "isLower":
			w.nl()
			w.println("func isLower(ch rune) bool {").indent()
			w.println("return ch >= 97 && ch <= 122").unindent()
			w.println("}")

		case "isPunct":
			w.nl()
			w.println("func isPunct(ch rune) bool {").indent()
			w.println("return (ch >= 33 && ch <= 47) || (ch >= 58 && ch <= 64) || (ch >= 91 && ch <= 96) || (ch >= 123 && ch <= 126)").unindent()
			w.println("}")

		case "isControl":
			w.nl()
			w.println("func isControl(ch rune) bool {").indent()
			w.println("return (ch >= 0 && ch <= 31) || ch == 127").unindent()
			w.println("}")

		case "isBlank":
			w.nl()
			w.println("func isBlank(ch rune) bool {").indent()
			w.println("return ch == 32 || ch == 9").unindent()
			w.println("}")
		case "isSpace":
			w.nl()
			w.println("func isSpace(ch rune) bool {").indent()
			w.println("return ch == 32 || (ch >= 9 && ch <= 13)").unindent()
			w.println("}")

		case "isGraph":
			w.nl()
			w.println("func isGraph(ch rune) bool {").indent()
			w.println("return isAlNum(ch) || isPunct(ch)").unindent()
			w.println("}")
			g.functions = AddUnique(g.functions, "isAlNum")
			g.functions = AddUnique(g.functions, "isPunct")

		case "isXDigit":
			w.nl()
			w.println("func isXDigit(ch rune) bool {").indent()
			w.println("return (ch >= 48 && ch <= 57) || (ch >= 65 && ch <= 70) || (ch >= 97 && ch <= 102)").unindent()
			w.println("}")
		}
	}
}

func (g *GoGenerator) DoGenerateProlog(vars any) {
	const tpl string = `
func Scanner(s string) error {
	oldstate := ""

	d := ScannerData{
		state: "0",
		ch: 0,
		token: "",
		cnt: 0,
	}

	runes := []rune(s)
	runes = append(runes, 0)

	for p := 0;p < len(runes); p++ {
		if d.state == oldstate {
			d.cnt ++
		} else {
			d.cnt = 0
		}

		oldstate = d.state
		d.ch = runes[p]

		switch d.state {
`
	g.writer.printTemplate(tpl, vars)
	g.writer.indentation = 4
}

func (g *GoGenerator) DoGenerateEpilog(finalState string) {
	w := &(g.writer)
	w.unindent()
	w.println("} else {").indent()
	w.println("return fmt.Errorf(\"Unexpected Char\")").unindent()
	w.println("} // if ").unindent().unindent()
	w.println("} // switch ").unindent()
	w.println("} // for")
	w.println("if d.state != " + g.Quote(finalState) + " {").indent()
	w.println("return fmt.Errorf(\"Unexpected End\")").unindent()
	w.println("}")
	w.println("return nil").unindent()
	w.println("}").unindent()
}

func (g *GoGenerator) DoClosePreviousIf() {
	w := &(g.writer)
	w.unindent()
	w.println("} else {").indent()
	w.println("return fmt.Errorf(\"Unexpected Char\")").unindent()
	w.println("}").nl()
}

func (g *GoGenerator) DoStartNewState(state string) {
	w := &(g.writer)
	w.unindent()
	w.println("case " + g.Quote(state) + ":")
	w.indent()
	w.print("if ")
}

func (g *GoGenerator) DoElseIf() {
	w := &(g.writer)
	w.unindent()
	w.print("} else if ")
}

func (g *GoGenerator) DoTestCharset(charset string) {
	fn := ""
	charset = strings.ToUpper(charset)
	switch charset {
	case "ALPHA":
		fn = "isAlpha"
	case "ALNUM":
		fn = "isAlNum"
	case "BLANK":
		fn = "isBlank"
	case "CONTROL":
		fn = "isControl"
	case "DIGIT":
		fn = "isDigit"
	case "GRAPH":
		fn = "isGraph"
	case "LOWER":
		fn = "isLower"
	case "PRINT":
		fn = "isPrint"
	case "PUNCT":
		fn = "isPunct"
	case "UPPER":
		fn = "isUpper"
	case "XDIGIT":
		fn = "isXDigit"
	}
	if fn != "" {
		g.writer.print(fn)
		g.writer.print("(d.ch)")
		g.AddWithDependencies(fn)
	}
}

func (g *GoGenerator) DoTestAny() {
	g.writer.println("true {").indent()
}

func (g *GoGenerator) DoMaxLoops(maxloops int) {
	w := &(g.writer)
	w.println(fmt.Sprintf("if d.cnt > %d {", maxloops)).indent()
	w.println("return fmt.Errorf(\"Unexpected Char\")").unindent()
	w.println("}")
}

func (g *GoGenerator) DoAction(action string, transition string, useLoops bool, testMode bool) {
	if testMode {
		g.writer.println(fmt.Sprintf("fmt.Println(\"Appel de '%s' pour la transition '%s'\")", g.Escape(action), g.Escape(transition)))
	} else {
		g.writer.println(fmt.Sprintf("OnAction(%s, %s, &d)", g.Quote(action), g.Quote(transition)))
	}
}

func (g *GoGenerator) DoResetToken() {
	g.writer.println("d.token = \"\"")
}

func (g *GoGenerator) DoAddToToken() {
	g.writer.println("d.token = d.token + string(d.ch)")
}

func (g *GoGenerator) DoNewState(state string) {
	g.writer.println("d.state = " + g.Quote(state))
}

func (g *GoGenerator) DoPrototype(actions []string) {
	w := &(g.writer)
	w.println("func OnAction(sTransition string, sAction string, data *ScannerData) {").indent()
	w.println("switch sAction {").indent()
	for _, action := range actions {
		w.println("case " + g.Quote(action) + ":").indent()
		w.println("// votre code ici").nl().unindent()
	}
	w.unindent()
	w.println("}").unindent()
	w.println("}")
}

func (g *GoGenerator) VisitCompare(c CompareInterface) {
	w := &(g.writer)
	switch c := c.(type) {
	case *CompareEos:
		w.print("int(d.ch) == 0")
	case CompareAny:
		w.print("true")
	case *CompareChar:
		w.print("d.ch == '" + c.ch + "'")
	case *CompareCharset:
		g.DoTestCharset(c.Name)
	case *CompareIn:
		w.print("strings.ContainsRune(" + g.Quote(c.Expression) + ", d.ch)")
	case *CompareOr:
		w.print("(")
		g.testLevel++
		for index, child := range c.childs {
			if index > 0 {
				w.print(" || ")
			}
			child.accept(g)
		}
		g.testLevel--
		w.print(")")
	}

	if g.testLevel == 0 {
		w.println(" {").indent()
	}
}

func (g *GoGenerator) AddWithDependencies(fn string) {
	switch fn {
	case "isAlpha":
		g.functions = AddUnique(g.functions, fn)
		g.AddWithDependencies("isUpper")
		g.AddWithDependencies("isLower")

	case "isAlNum":
		g.functions = AddUnique(g.functions, fn)
		g.AddWithDependencies("isAlpha")
		g.AddWithDependencies("isDigit")

	case "isDigit":
		g.functions = AddUnique(g.functions, fn)

	case "isUpper":
		g.functions = AddUnique(g.functions, fn)

	case "isLower":
		g.functions = AddUnique(g.functions, fn)

	case "isPunct":
		g.functions = AddUnique(g.functions, fn)

	case "isControl":
		g.functions = AddUnique(g.functions, fn)

	case "isBlank":
		g.functions = AddUnique(g.functions, fn)
	case "isSpace":
		g.functions = AddUnique(g.functions, fn)

	case "isGraph":
		g.functions = AddUnique(g.functions, fn)
		g.AddWithDependencies("isAlNum")
		g.AddWithDependencies("isPunct")

	case "isXDigit":
		g.functions = AddUnique(g.functions, fn)
	}
}
