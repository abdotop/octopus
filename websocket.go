package octopus

// octopus "real_time_forum/app"

func (app *App) WS(path string, handler ...HandlerFunc) {
	app.handle(path, handler, "GET")
}
