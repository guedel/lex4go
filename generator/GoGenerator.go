package generator

import (
	"fmt"
	"strings"
)

type GoGenerator struct {
	writer    CodeWriter
	testLevel int
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
	"regexp"
)	

type ScannerData struct {
	state string
	ch string
	token string
	cnt int
}

func like(ch string, expr string) bool {
	result,_ := regexp.MatchString(expr, ch)
	return result
}
`
	g.writer.printTemplate(tpl, vars)
	g.writer.indentation = 1
}

func (g *GoGenerator) DoEndDocument(vars any) {
}

func (g *GoGenerator) DoGenerateProlog(vars any) {
	const tpl string = `
func Scanner(s string) error {
	oldstate := ""

	d := ScannerData{
		state: "0",
		ch: "",
		token: "",
		cnt: 0,
	}

	s = s + string(0)

	for p := 0;p < len(s); p++ {
		if d.state == oldstate {
			d.cnt ++
		} else {
			d.cnt = 0
		}

		oldstate = d.state
		d.ch = s[p]

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
	var test string
	switch strings.ToUpper(charset) {
	case "ALPHA":
		test = "isAlpha(ch)"
	case "BLANK":
		test = "isBlank(ch)"
	case "CONTROL":
		test = "isControl(ch)"
	case "DIGIT":
		test = "isDigit(ch)"
	case "GRAPH":
		test = "isGraph(ch)"
	case "LOWER":
		test = "isLower(ch)"
	case "PRINT":
		test = "isPrint(ch)"
	case "PUNCT":
		test = "isPunct(ch)"
	case "UPPER":
		test = "isUpper(ch)"
	case "XDIGIT":
		test = "isXDigit(ch)"
	}
	g.writer.print(test)
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
	g.writer.println("d.token = d.token + d.ch")
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
		w.print("d.ch dans (" + c.Expression + ")")
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
		/*
			case In:
				w.print("Dans (")
				w.print(")")
		*/
	}

	if g.testLevel == 0 {
		w.println(" {").indent()
	}
}
