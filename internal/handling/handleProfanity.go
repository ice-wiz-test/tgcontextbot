package handling

import (
	"strings"
)

func CheckProf(badWords *[]string, find string) bool {
	for i := 0; i < len(*(badWords)); i++ {
		if strings.Contains((*badWords)[i], find) {
			return true
		}
	}
	return false
}
