package log

import (
	"github.com/skoo87/log4go"
)

var (
	Logger  = log4go.NewLogger()
	console = log4go.NewConsoleWriter()
	file    = log4go.NewFileWriter()
	name    = "cacuserver"
	path    = "./logs"
)

func SetLevel(level string) {
	l := log4go.DEBUG
	switch level {
	case "debug":
		l = log4go.DEBUG
	case "info":
		l = log4go.INFO
	case "warn":
		l = log4go.WARNING
	case "error":
		l = log4go.ERROR
	case "fatal":
		l = log4go.FATAL
	}
	Logger.SetLevel(l)
}

func SetLogPath(p, n string) {
	if p[len(p)-1] != '/' || p[len(p)-1] != '\\' {
		p = p + "/"
	}
	path = p
	name = n
}

func Run() {
	file.SetPathPattern(path + name + "-%Y-%M-%D-%H.log")
	console.SetColor(true)
	Logger.Register(console)
	Logger.Register(file)
	Logger.SetLayout("2006-01-02 15:04:05")
}
