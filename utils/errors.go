package utils

import "log"

func Fail(msg string) {
	log.Fatal("ERROR: ", msg)
}
func FailIf(err error, msg string) {
	if err != nil {
		log.Fatal("ERROR: ", msg, "(", err, ")")
	}
}