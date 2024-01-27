package helpers

import "fmt"

func IsClientKey(s string) bool {
	if len(s) != 64 {
		return false
	}
	for _, char := range s {
		if char < 48 || (char > 58 && char < 97) || char > 122 {
			fmt.Println(char)
			return false
		}
	}
	return true
}
