package logger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	debug  = log.New(ioutil.Discard, "INF: ", log.LstdFlags|log.Lshortfile)
	info   = log.New(os.Stderr, "INF: ", log.LstdFlags|log.Lshortfile)
	errlog = log.New(os.Stderr, "ERR: ", log.LstdFlags|log.Lshortfile)
)

// GetErrLogger ...
func GetErrLogger() *log.Logger {
	return errlog
}

func Debug(v ...interface{}) {
	debug.Output(2, fmt.Sprint(v...))
}
func Debugf(format string, v ...interface{}) {
	debug.Output(2, fmt.Sprintf(format, v...))
}
func Debugln(v ...interface{}) {
	debug.Output(2, fmt.Sprintln(v...))
}
func Print(v ...interface{}) {
	info.Output(2, fmt.Sprint(v...))
}
func Printf(format string, v ...interface{}) {
	info.Output(2, fmt.Sprintf(format, v...))
}
func Println(v ...interface{}) {
	info.Output(2, fmt.Sprintln(v...))
}
func Error(v ...interface{}) {
	errlog.Output(2, fmt.Sprint(v...))
}
func Errorf(format string, v ...interface{}) {
	errlog.Output(2, fmt.Sprintf(format, v...))
}
func Errorln(v ...interface{}) {
	errlog.Output(2, fmt.Sprintln(v...))
}
func Fatal(v ...interface{}) {
	errlog.Output(2, fmt.Sprint(v...))
}
func Fatalf(format string, v ...interface{}) {
	errlog.Output(2, fmt.Sprintf(format, v...))
}
func Fatalln(v ...interface{}) {
	errlog.Output(2, fmt.Sprintln(v...))
}
