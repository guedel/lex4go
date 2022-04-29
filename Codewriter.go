package main

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
}

func newCodeWriter() CodeWriter {
	return CodeWriter{
		spacer:      " ",
		indentation: 0,
		current:     "",
	}
}

func (c *CodeWriter) indent() *CodeWriter {
	c.indentation++
	return c
}

func (c *CodeWriter) unindent() *CodeWriter {
	c.indentation--
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
	err := t.Execute(os.Stdout, vars)
	if err != nil {
		log.Println("template error: ", err)
	}
	return c
}

func (c *CodeWriter) nl() *CodeWriter {
	fmt.Print(strings.Repeat(c.spacer, int(c.indentation)))
	fmt.Println(c.current)
	c.current = ""
	return c
}
