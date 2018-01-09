package util

import "log"

func Check(e error) (bool) {
	if e != nil {
		log.Fatalln(e)
		return false
	}
	return true
}