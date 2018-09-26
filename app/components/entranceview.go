package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"

	"github.com/nobonobo/gophertank/app/actions"
	"github.com/nobonobo/gophertank/app/dispatcher"
)

// EntranceView ...
type EntranceView struct {
	vecty.Core
	Name string `vecty:"prop"`
}

// Render ...
func (c *EntranceView) Render() vecty.ComponentOrHTML {
	vecty.SetTitle("Join")
	return elem.Body(
		vecty.Markup(
			prop.ID("join-form-body"),
			vecty.Class("h-100", "text-center", "d-flex", "justify-content-center", "align-content-center"),
		),
		elem.Form(
			vecty.Markup(
				event.Submit(func(ev *vecty.Event) {
					ev.Call("preventDefault")
					name := ev.Target.Get("displayName").Get("value").String()
					dispatcher.Dispatch(&actions.Enter{Name: name})
				}),
			),
			elem.Image(
				vecty.Markup(
					vecty.Class("mb-4"),
					prop.Src("favicon.png"),
					vecty.Attribute("alt", ""),
					vecty.Attribute("width", "72"),
					vecty.Attribute("height", "72"),
				),
			),
			elem.Heading1(
				vecty.Markup(
					vecty.ClassMap{
						"h3":   true,
						"mb-3": true,
					},
				),
				vecty.Text("Gopher Tank"),
			),
			elem.Heading5(
				vecty.Text("Multiplayer Real Time Online Game"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("form-group"),
				),
				elem.Label(
					vecty.Markup(
						prop.For("displayName"),
					),
					vecty.Text("Display Name"),
				),
				elem.Input(
					vecty.Markup(
						prop.Type("text"),
						prop.ID("displayName"),
						vecty.Class("form-control"),
						prop.Placeholder("Display Name"),
						vecty.Attribute("required", ""),
						vecty.Attribute("minlength", "1"),
						vecty.Attribute("maxlength", "16"),
						prop.Autofocus(true),
						prop.Value(c.Name),
					),
				),
			),
			elem.Button(
				vecty.Markup(
					vecty.ClassMap{
						"btn":         true,
						"btn-lg":      true,
						"btn-primary": true,
						"btn-block":   true,
					},
					prop.Type("submit"),
				),
				vecty.Text("Enter"),
			),
			elem.Paragraph(
				vecty.Markup(
					vecty.ClassMap{
						"mt-5":       true,
						"mb-3":       true,
						"text-muted": true,
					},
				),
				vecty.Text("© 2018 nobonobo"),
			),
		),
	)
}
