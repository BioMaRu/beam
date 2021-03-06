package hime

// Globals is the global const map
type Globals map[interface{}]interface{}

func cloneGlobals(xs Globals) Globals {
	if xs == nil {
		return nil
	}
	rs := make(Globals)
	for k, v := range xs {
		rs[k] = v
	}
	return rs
}

// Globals registers global constants
func (app *App) Globals(globals Globals) *App {
	if app.globals == nil {
		app.globals = make(Globals)
	}
	for key, value := range globals {
		app.globals[key] = value
	}
	return app
}

// Global gets value from global storage
func (app *App) Global(key interface{}) interface{} {
	if app.globals == nil {
		return nil
	}
	return app.globals[key]
}

// Global returns global value
func (ctx *Context) Global(key interface{}) interface{} {
	return ctx.app.Global(key)
}
