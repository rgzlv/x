package util

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var errLogger log.Logger = *log.New(os.Stderr, "[ERROR] ", log.Ltime|log.Lshortfile)
var infoLogger log.Logger = *log.New(os.Stderr, "[INFO] ", log.Ltime|log.Lshortfile)

func LogError(err string) {
	errLogger.Output(2, err)
}

func LogErrorf(format string, a ...any) {
	errLogger.Output(2, fmt.Sprintf(format, a...))
}

func LogErrorDepth(depth int, err string) {
	errLogger.Output(depth, err)
}

func LogHTTPError(w http.ResponseWriter, err error) {
	errLogger.Output(2, err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func LogFatal(err string) {
	errLogger.Output(2, err)
	os.Exit(1)
}

func LogInfo(a ...any) {
	infoLogger.Output(2, fmt.Sprintln(a...))
}

func LogInfof(format string, a ...any) {
	infoLogger.Output(2, fmt.Sprintf(format, a...))
}
