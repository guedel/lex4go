package main

import "github.com/guedel/lex4go/generator"

type LanguageType uint

type GeneratorError struct {
	message string
}

func (e GeneratorError) Error() string {
	return e.message
}

const (
	Algorithm = iota
	VisualBasic
	Php
	Go
	// TODO
)

func (l LanguageType) getGenerator() generator.GeneratorInterface {
	switch l {
	case Algorithm:
		return &generator.AlgorithmGenerator{}
	}
	return nil
}

/**
 * Procédure de génération de la machine à états.
 */
func GenerateStateEngine(lexer Lexer, language LanguageType, testMode bool, genProto bool) error {
	gen := language.getGenerator()

	var vars map[string]string
	var useLoop bool

	finalState := ""

	var actions []string
	gen.DoStartDocument(lexer)
	gen.DoGenerateProlog(lexer)
	oldState := ""
	for _, rule := range lexer.Rules {
		state := rule.From
		transition := rule.Id
		if state != oldState {
			if oldState != "" {
				gen.DoClosePreviousIf()
			}
			gen.DoStartNewState(rule.From)
		} else {
			gen.DoElseIf()
		}
		oldState = state
		gen.VisitCompare(rule.Test.Compare)
		if rule.Repeat > 0 {
			useLoop = true
			gen.DoMaxLoops(rule.Repeat)
		} else {
			useLoop = false
		}

		if rule.Action != "" {
			s := rule.Action
			gen.DoAction(s, transition, useLoop, testMode)
			if genProto {
				actions = addUnique(actions, rule.Action)
			}
		}

		if rule.Reset {
			gen.DoResetToken()
		}
		if rule.Concat {
			gen.DoAddToToken()
		}
		if rule.To != state {
			gen.DoNewState(rule.To)
		}
		if rule.Final {
			if len(finalState) > 0 {
				return GeneratorError{"Duplicate final state"}
			}
			finalState = rule.From
		}
	}
	gen.DoGenerateEpilog(finalState)
	if genProto {
		gen.DoPrototype(actions)
	}
	gen.DoEndDocument(vars)
	return nil
}

/**
 * N'ajoute que si le token n'existe pas encore
 */
func addUnique(vars []string, token string) []string {
	find := false
	for _, value := range vars {
		if value == token {
			find = true
		}
	}
	if !find {
		return append(vars, token)
	}
	return vars
}
