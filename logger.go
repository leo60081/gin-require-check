package middleware

import "log"
import "os"

var logger Logger

func init() {
	logger = log.New(os.Stderr, "[GIN-ReqCheck] ", log.LstdFlags)
}

func SetLogger(l Logger) {
	logger = l
	return
}

type Logger interface {
	Print(v ...interface{})
	Println(v ...interface{})
}
