package handling

import (
	"strings"
)

func CheckProf(badWords *[]string, find string, exceptedWords *[]string) bool {
	for i := 0; i < len(*(badWords)); i++ {
		for j := 0; j < len(*exceptedWords); j++ {
			if (*badWords)[i] == (*exceptedWords)[j] {
				continue
			}
			if strings.Contains(find, (*badWords)[i]) {
				return true
			}
		}
	}
	return false
}
