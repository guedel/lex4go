package main

import "encoding/xml"

type Lexer struct {
	XMLName      xml.Name `xml:"lexer"`
	Name         string   `xml:"name"`
	Author       string   `xml:"author"`
	Description  string   `xml:"description"`
	CreationDate string   `xml:"dateCreation"`
	UpdateDate   string
	Initial      string `xml:"initial"`
	Rules        []Rule `xml:"rules>rule"`
}

/*
func (l *Lexer) UnmarshalText(text []byte) error {
	fmt.Println("Lexer.UnmarshalText:", string(text))
	return nil
}
*/
