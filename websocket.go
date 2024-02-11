package octopus

// octopus "real_time_forum/App"

func (a *App) WS(path string, handler ...Handler) {
	a.handle(path, handler, "GET")
}
