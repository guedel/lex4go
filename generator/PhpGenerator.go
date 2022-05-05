package generator

import (
	"fmt"
	"strings"
)

type PhpGenerator struct {
	writer    CodeWriter
	testLevel int
}

func (g *PhpGenerator) Quote(c string) string {
	return "\"" + g.Escape(c) + "\""
}

func (g *PhpGenerator) Escape(c string) string {
	return strings.ReplaceAll(c, "\"", "\\\"")

}

func (g *PhpGenerator) DoStartDocument(vars any) {
	if g.writer == (CodeWriter{}) {
		g.writer = newCodeWriter()
		g.writer.spacer = "\t"
	}
	const tpl string = `
<?php	
	// {{.Name}}	
	// ----------------------------
	// Création de {{.Author}}
	// Créé le {{.CreationDate}}
	// Modifié le {{.UpdateDate}}
	// ----------------------------
	// {{.Description}}

	class ScannerData {
		public string $state = "{{.Initial}}";
		public string $ch = '';
		public string $token = "";
		public int $cnt = 0;
	}

	class UnexpectedCharException extends \Exception { }
	class UnexpectedEndException extends \Exception { }
`
	g.writer.printTemplate(tpl, vars)
	g.writer.indentation = 2
}

func (g *PhpGenerator) DoEndDocument(vars any) {
}

func (g *PhpGenerator) DoGenerateProlog(vars any) {
	const tpl string = `
	function Scanner(string $s) 
	{
		try {
			$oldstate = null;
			$d = new ScannerData();
			$s .= chr(0);

			for ($p = 0; $p < strlen($s); $p++) {
				if ($d->state === $oldstate) {
					$d->cnt++;
				} else {
					$d->cnt = 0;
				}
				$oldstate = $d->state;
				$d->ch = substr($s, $p, 1);

				switch ($d->state) {
`
	g.writer.printTemplate(tpl, vars)
	g.writer.indentation = 6
}

func (g *PhpGenerator) DoGenerateEpilog(finalState string) {
	w := &(g.writer)
	w.unindent()
	w.println("} else {").indent()
	w.println("throw new UnexpectedCharException();").unindent()
	w.println("}").println("break;").unindent().unindent()
	w.println("}").unindent()
	w.println("}")
	w.println("if ($d->state !== " + g.Quote(finalState) + ") {").indent()
	w.println("throw new UnexpectedEndException();").unindent()
	w.println("}").unindent()
	w.println("} catch (UnexpectedEndException $ex) {").indent()
	w.println("throw new \\Exception('Fin inattendue');").unindent()
	w.println("} catch (UnexpectedCharException $ex) {").indent()
	w.println("throw new \\Exception(\"Caractère '$d->ch' inattendu en position $p\");").unindent()
	w.println("}").unindent()
	w.println("}").nl()

	/*
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
	*/
}

func (g *PhpGenerator) DoClosePreviousIf() {
	w := &(g.writer)
	w.unindent()
	w.println("} else {").indent()
	w.println("throw new UnexpectedCharException();").unindent()
	w.println("}")
	w.println("break;").nl()
}

func (g *PhpGenerator) DoStartNewState(state string) {
	w := &(g.writer)
	w.unindent()
	w.println("case " + g.Quote(state) + ":")
	w.indent()
	w.print("if ")
}

func (g *PhpGenerator) DoElseIf() {
	w := &(g.writer)
	w.unindent()
	w.print("} elseif ")
}

func (g *PhpGenerator) DoTestCharset(charset string) {
	var test string
	switch strings.ToUpper(charset) {
	case "ALPHA":
		test = "isAlpha($d->ch)"
	case "BLANK":
		test = "isBlank($d->ch)"
	case "CONTROL":
		test = "isControl($d->ch)"
	case "DIGIT":
		test = "isDigit($d->ch)"
	case "GRAPH":
		test = "isGraph($d->ch)"
	case "LOWER":
		test = "isLower($d->ch)"
	case "PRINT":
		test = "isPrint($d->ch)"
	case "PUNCT":
		test = "isPunct($d->ch)"
	case "UPPER":
		test = "isUpper($d->ch)"
	case "XDIGIT":
		test = "isXDigit($d->ch)"
	}
	g.writer.print(test)
}

func (g *PhpGenerator) DoTestAny() {
	g.writer.println("(true) {").indent()
}

func (g *PhpGenerator) DoMaxLoops(maxloops int) {
	w := &(g.writer)
	w.println(fmt.Sprintf("if ($d->cnt > %d) {", maxloops)).indent()
	w.println("throw new UnexpectedCharException();").unindent()
	w.println("}")
}

func (g *PhpGenerator) DoAction(action string, transition string, useLoops bool, testMode bool) {
	if testMode {
		g.writer.println(fmt.Sprintf("echo 'call of  \\'%s\\' for transition \\'%s\\'';", g.Escape(action), g.Escape(transition)))
	} else {
		g.writer.println(fmt.Sprintf("%s(%s, $d);", action, g.Quote(transition)))
	}
}

func (g *PhpGenerator) DoResetToken() {
	g.writer.println("$d->token = \"\";")
}

func (g *PhpGenerator) DoAddToToken() {
	g.writer.println("$d->token .= $d->ch;")
}

func (g *PhpGenerator) DoNewState(state string) {
	g.writer.println("$d->state = " + g.Quote(state) + ";")
}

func (g *PhpGenerator) DoPrototype(actions []string) {
	w := &(g.writer)
	for _, action := range actions {
		w.println("function " + action + "($transition, ScannerData $d)")
		w.println("{").indent()
		w.println("// votre code ici").nl().unindent()
		w.println("}").nl()
	}
}

func (g *PhpGenerator) VisitCompare(c CompareInterface) {
	w := &(g.writer)
	w.print("(")
	switch c := c.(type) {
	case *CompareEos:
		w.print("ord($d->ch) === 0")
	case CompareAny:
		w.print("true")
	case *CompareChar:
		w.print("$d->ch === " + g.Quote(c.ch))
	case *CompareCharset:
		g.DoTestCharset(c.Name)
	case *CompareIn:
		w.print("strpos(" + g.Quote(c.Expression) + ", $d->ch)!==false")
	case *CompareOr:
		if g.testLevel > 0 {
			w.print("(")
		}
		g.testLevel++
		for index, child := range c.childs {
			if index > 0 {
				w.print(" || ")
			}
			child.accept(g)
		}
		g.testLevel--
		if g.testLevel > 0 {
			w.print(")")
		}
	}

	if g.testLevel == 0 {
		w.println(") {").indent()
	}
}
