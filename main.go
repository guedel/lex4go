package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/guedel/lex4go/generator"
	"log"
	"os"
	"strings"
)

/*
	Point d'entrée du programme.
	<nomprog> [options] fichier

	où options:
	- l, language : langage de génération. Par défaut algorithme.
	- h, help: affichage de l'aide
	- t: utilisation du mode test
	- p: génération du prototype de méthode pour les actions

*/

func main() {
	var language string
	var useHelp bool
	var testMode bool
	var genProto bool
	flag.StringVar(&language, "language", "algorithm", "Language to use")
	flag.StringVar(&language, "l", "algorithm", "Language to use (shorthand)")
	flag.BoolVar(&useHelp, "help", false, "Display help")
	flag.BoolVar(&useHelp, "h", false, "Display help (shorthand)")
	flag.BoolVar(&testMode, "t", false, "Generate with test mode")
	flag.BoolVar(&genProto, "p", false, "Generate with action prototype method")

	flag.Parse()
	if flag.Arg(0) == "" || useHelp {
		flag.PrintDefaults()
		return
	}

	var gen generator.GeneratorInterface
	switch strings.ToLower(language) {
	case "algorithm", "algo":
		gen = &generator.AlgorithmGenerator{}
	/*
		case "visualbasic", "vb":
		case "go":
		case "php":
	*/
	default:
		log.Fatal("Unknown language")
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

		GenerateStateEngine(lexer, gen, testMode, genProto)
	}
}
