package handling

import (
	"strings"
)

func CheckProf(badWords *[]string, find string) bool {
	for i := 0; i < len(*(badWords)); i++ {
		if strings.Contains(find, (*badWords)[i]) {
			return true
		}
	}
	return false
}
