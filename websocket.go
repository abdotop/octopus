package octopus

// octopus "real_time_forum/App"

func (a *App) WS(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "GET")
}
