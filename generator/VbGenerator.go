package generator

import (
	"fmt"
	"strings"
)

type VbGenerator struct {
	writer    CodeWriter
	testLevel int
}

func (g *VbGenerator) Quote(c string) string {
	return "\"" + g.Escape(c) + "\""
}

func (g *VbGenerator) Escape(c string) string {
	return strings.ReplaceAll(c, "\"", "\\\"")

}

func (g *VbGenerator) DoStartDocument(vars any) {
	if g.writer == (CodeWriter{}) {
		g.writer = newCodeWriter()
		g.writer.spacer = "\t"
	}
	const tpl string = `
' {{.Name}}	
' ----------------------------
' Création de {{.Author}}
' Créé le {{.CreationDate}}
' Modifié le {{.UpdateDate}}
' ----------------------------
' {{.Description}}

Option Compare Database
Option Explicit

Const Err_UnexpectedChar = 1 + vbObjectError + 512
Const Err_UnexpectedEnd = 2 + vbObjectError + 512

Type DataScanner
	var state: String
	var ch: Caractère
	var token: String
	var cnt: Integer
End Type
`
	g.writer.printTemplate(tpl, vars)
	g.writer.indentation = 2
}

func (g *VbGenerator) DoEndDocument(vars any) {
}

func (g *VbGenerator) DoGenerateProlog(vars any) {
	const tpl string = `
Sub Scanner(s as String)
	Dim old_state as String: old_state = ""
	Dim p As Integer: p = 0
	Dim d As DataScanner

	On Error Goto trap

	With d
		.state = "0"
		.ch = ''
		.token = ""
		.cnt = 0
	End With

	s = s & chr(0)

	Do While p <= Len(s)
		If d.state = old_state Then
			d.cnt = d.cnt + 1
		Else
			d.cnt = 0
		End If

		old_state = d.state
		d.ch = Mid$(s, p, 1)

		Select Case d.state
`
	g.writer.printTemplate(tpl, vars)
	g.writer.indentation = 4
}

func (g *VbGenerator) DoGenerateEpilog(finalState string) {
	w := &(g.writer)
	w.unindent()
	w.println("Else").indent()
	w.println("Err.Raise Exception_UnexpectedChar").unindent()
	w.println("End If").unindent()
	w.println("End Select")
	w.println("p = p + 1").unindent()
	w.println("Loop")
	w.println("If d.state <> " + g.Quote(finalState) + " Then").indent()
	w.println("Err.Raise Err_UnexpectedEnd")
	w.println("Exit Sub").unindent()
	w.println("End If").unindent()
	w.nl().unindent()
	w.println("trap:").indent()
	w.println("If Err.Number = Err_UnexpectedChar Or ch = vbNullChar Then").indent()
	w.println("Err.Descritpion = \"Fin de chaine inattendue\"").unindent()
	w.println("Else").indent()
	w.println("Err.Description = \"Caractère inattendu \"\" & ch & \" en position \" p").unindent()
	w.println("End If")
	w.println("Err.Raise Err.Number").unindent()
	w.println("End Sub")
}

func (g *VbGenerator) DoClosePreviousIf() {
	w := &(g.writer)
	w.unindent()
	w.println("Else").indent()
	w.println("Err.Raise Err_UnexcpectedChar").unindent()
	w.println("End If").nl()
}

func (g *VbGenerator) DoStartNewState(state string) {
	w := &(g.writer)
	w.unindent()
	w.println("Case " + g.Quote(state) + ":")
	w.indent()
	w.print("If ")
}

func (g *VbGenerator) DoElseIf() {
	w := &(g.writer)
	w.unindent()
	w.print("Else If ")
}

func (g *VbGenerator) DoTestLike(expr string) {
	g.writer.println("d.ch Like " + g.Quote(expr) + " Then").indent()
}

func (g *VbGenerator) DoTestCharset(charset string) {
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
		test = "isPunt(ch)"
	case "UPPER":
		test = "isUpper(ch)"
	case "XDIGIT":
		test = "isXDigit(ch)"
	}
	g.writer.print(test)
}

func (g *VbGenerator) DoTestAny() {
	g.writer.println("true Then").indent()
}

func (g *VbGenerator) DoMaxLoops(maxloops int) {
	w := &(g.writer)
	w.println(fmt.Sprintf("If d.cnt > %d Then", maxloops)).indent()
	w.println("Err.Raise Err_UnexpectedChar").unindent()
	w.println("End If")
}

func (g *VbGenerator) DoAction(action string, transition string, useLoops bool, testMode bool) {
	if testMode {
		g.writer.println(fmt.Sprintf("Print \"Appel de '%s' pour la transition '%s'\"", g.Escape(action), g.Escape(transition)))
	} else {
		g.writer.println(fmt.Sprintf("OnAction %s, %s, d", g.Quote(action), g.Quote(transition)))
	}
}

func (g *VbGenerator) DoResetToken() {
	g.writer.println("d.token = \"\"")
}

func (g *VbGenerator) DoAddToToken() {
	g.writer.println("d.token = d.token & d.ch")
}

func (g *VbGenerator) DoNewState(state string) {
	g.writer.println("d.state = " + g.Quote(state))
}

func (g *VbGenerator) DoPrototype(actions []string) {
	w := &(g.writer)
	w.println("Sub OnAction(sTransition as String, sAction as String, var data As DataScanner)").indent()
	w.println("Select Case sAction").indent()
	for _, action := range actions {
		w.println("Case " + g.Quote(action)).indent()
		w.println("' votre code ici").nl().unindent()
	}
	w.unindent()
	w.println("End Select").unindent()
	w.println("End Sub")
}

func (g *VbGenerator) VisitCompare(c CompareInterface) {
	w := &(g.writer)
	switch c := c.(type) {
	case *CompareEos:
		w.print("Asc(d.ch) = 0")
	case CompareAny:
		w.print("True")
	case *CompareChar:
		w.print("d.ch = " + g.Quote("ch"))
	case *CompareCharset:
		g.DoTestCharset(c.Name)
	case *CompareLike:
		w.print("d.ch Like " + g.Quote(c.Expression))
	case *CompareAnd:
		w.print("(")
		g.testLevel++
		for index, child := range c.childs {
			if index > 0 {
				w.print(" And ")
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
				w.print(" Or ")
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
		w.println(" Then").indent()
	}
}
