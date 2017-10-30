package emulator

import (
	"github.com/QSpike/cmd"
)

type Application struct {
	agents map[string]*Agent
	closed bool
}

func New() *Application {
	app := &Application{
		agents: make(map[string]*Agent),
		closed: false,
	}
	app.init()
	return app
}

func (app *Application) init() {
}

func (app *Application) Close() {
	if app.closed {
		return
	}
}

func (app *Application) Run() {
	app.mainLoop()
}

func (app *Application) mainLoop() {
	c := cmd.New(new(Client))
	c.Cmdloop()
}
