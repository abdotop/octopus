package octopus

// octopus "real_time_forum/app"

func (a *app) WS(path string, handler ...HandlerFunc) {
	a.handle(path, handler, "GET")
}
