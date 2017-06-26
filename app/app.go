package app

import (
	"github.com/kolonse/CacuServer/log"
)

type module interface {
	Init()
	Run() error
	Stop()
	Exit(error)
}

type App struct {
	group []module
}

func (app *App) Init(m ...module) {
	app.group = append(app.group, m...)
}

func (app *App) Go() {
	c := make(chan int, 1)
	ok := make(chan int, 1)
	for i := 0; i < len(app.group); i++ {
		app.group[i].Init()
		go func(index int) {
			ok <- index
			app.group[index].Exit(app.group[index].Run())
			c <- index
		}(i)
	}
	for i := 0; i < len(app.group); i++ {
		<-ok
	}
	log.Logger.Info("应用启动完成")
	for i := 0; i < len(app.group); i++ {
		<-c
	}
	log.Logger.Info("应用退出完成")
}

func (app *App) Stop() {
	for i := 0; i < len(app.group); i++ {
		app.group[i].Stop()
	}
}

var GlobalApp App

func Init(m ...module) {
	GlobalApp.Init(m...)
}

func Go() {
	GlobalApp.Go()
}

func Stop() {
	GlobalApp.Stop()
}
