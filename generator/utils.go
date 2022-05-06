package generator

/**
 * N'ajoute que si le token n'existe pas encore
 */
func AddUnique(vars []string, token string) []string {
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
