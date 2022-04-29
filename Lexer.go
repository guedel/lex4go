package main

type Lexer struct {
	Name         string `xml:"name"`
	Author       string `xml:"author"`
	Description  string `xml:"description"`
	CreationDate string
	UpdateDate   string
	Rules        []Rule `xml:"rule"`
}

/*
func (l *Lexer) UnmarshalText(text []byte) error {
	fmt.Println("Lexer.UnmarshalText:", string(text))
	return nil
}
*/
