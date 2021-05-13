package logging

import "log"

const INFO = "INFO"
const DEBUG = "DEBUG"
const ERROR = "ERROR"

func Log(lvl string, msg string) {
	log.Printf("%s %s", lvl, msg)
}