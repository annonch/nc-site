package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// hello is a component that displays a simple "Hello World!". A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type hello struct {
	app.Compo
}

type myCompo struct {
	app.Compo
}

// The Render method is where the component appearance is defined. Here, a
// "Hello World!" is displayed as a heading.
// func (h *hello) Render() app.UI {
// 	return app.H1().Text("Hello World!")
// }

func (c *myCompo) Render() app.UI {
	return app.Div().Body(
		app.H1().
			Class("title").
			Text("Build a GUI with Go"),
		app.P().
			Class("text").
			Text("Just because Go and this package are really awesome!"),
	)
}
