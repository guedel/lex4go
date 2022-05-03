package generator

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

type CodeWriter struct {
	spacer      string
	indentation uint
	current     string
	out         os.File
}

func newCodeWriter() CodeWriter {
	return CodeWriter{
		spacer:      "    ",
		indentation: 0,
		current:     "",
		out:         *os.Stdout,
	}
}

func (c *CodeWriter) indent() *CodeWriter {
	c.indentation++
	return c
}

func (c *CodeWriter) unindent() *CodeWriter {
	if c.indentation > 0 {
		c.indentation--
	}
	return c
}

func (c *CodeWriter) print(msg string) *CodeWriter {
	c.current += msg
	return c
}

func (c *CodeWriter) println(msg string) *CodeWriter {
	c.current += msg
	c.nl()
	return c
}

func (c *CodeWriter) printTemplate(tpl string, vars any) *CodeWriter {
	t := template.Must(template.New("template").Parse(tpl))
	err := t.Execute(&c.out, vars)
	if err != nil {
		log.Println("template error: ", err)
	}
	return c
}

func (c *CodeWriter) nl() *CodeWriter {
	fmt.Fprint(&c.out, strings.Repeat(c.spacer, int(c.indentation)))
	fmt.Fprintln(&c.out, c.current)
	c.current = ""
	return c
}
