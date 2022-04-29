package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
	Point d'entrée du programme.
	<nomprog> [options] fichier

	où options:
	- l, language : langage de génération. Par défaut algorithme.

*/

func main() {
	var language string
	var useHelp bool
	var languageToUse LanguageType
	flag.StringVar(&language, "language", "algorithm", "Language to use")
	flag.StringVar(&language, "l", "algorithm", "Language to use (shorthand)")
	flag.BoolVar(&useHelp, "help", false, "Display help")
	flag.BoolVar(&useHelp, "h", false, "Display help (shorthand)")

	flag.Parse()
	if flag.Arg(0) == "" || useHelp {
		flag.PrintDefaults()
	}

	switch strings.ToLower(language) {
	case "algorithm", "algo":
		languageToUse = Algorithm
	case "visualbasic", "vb":
		languageToUse = VisualBasic
	case "go":
		languageToUse = Go
	case "php":
		languageToUse = Php
	}

	fmt.Printf("Language choice: %s\n", language)

	for _, filename := range flag.Args() {
		source, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		var lexer Lexer
		if err := xml.Unmarshal(source, &lexer); err != nil {
			log.Fatal(err)
		}
		// fmt.Println("Main source:", string(source))

		if languageToUse == Algorithm {

		}

		fmt.Printf("%v\n", lexer)
	}

	// TODO:
	// - lire le/s fichier/s
	// - lancer la génération
}
