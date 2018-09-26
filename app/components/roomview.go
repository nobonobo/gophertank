package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"

	"github.com/nobonobo/gophertank/app/actions"
	"github.com/nobonobo/gophertank/app/dispatcher"
	"github.com/nobonobo/gophertank/app/store"
)

// RoomView ...
type RoomView struct {
	vecty.Core
}

// Render ...
func (c *RoomView) Render() vecty.ComponentOrHTML {
	vecty.SetTitle("Room")
	items := vecty.List{}
	for _, item := range store.CurrentRoomMembers {
		items = append(items,
			elem.ListItem(
				vecty.Markup(
					vecty.ClassMap{
						"list-group-item":         true,
						"d-flex":                  true,
						"justify-content-between": true,
						"align-items-center":      true,
					},
				),
				vecty.Text(item.Name),
				elem.Span(
					vecty.Markup(
						vecty.ClassMap{
							"badge":         true,
							"badge-primary": true,
							"badge-pill":    true,
						},
					),
					vecty.Text(item.UUID),
				),
			),
		)
	}

	return elem.Body(
		vecty.Markup(
			vecty.Style("padding-top", "5rem"),
		),
		elem.Navigation(
			vecty.Markup(
				vecty.ClassMap{
					"navbar":           true,
					"navbar-expand-md": true,
					"navbar-dark":      true,
					"bg-dark":          true,
					"fixed-top":        true,
				},
			),
			vecty.If(
				len(store.CurrentRoomMembers) < 2,
				elem.Anchor(
					vecty.Markup(
						vecty.ClassMap{
							"navbar-brand": true,
						},
						prop.Href("#"),
						event.Click(func(ev *vecty.Event) {
							ev.Call("preventDefault")
							dispatcher.Dispatch(&actions.Leave{})
						}),
					),
					elem.Italic(
						vecty.Markup(
							vecty.Class("material-icons"),
						),
						vecty.Text("arrow_back"),
					),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.ClassMap{
						"mx-auto": true,
					},
				),
				elem.Heading1(
					vecty.Markup(
						vecty.Class("navbar-brand"),
					),
					vecty.Text("Room"),
				),
			),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("container"),
			),
			elem.Heading2(
				vecty.If(store.Wanted, vecty.Text("Wanted Member!")),
				vecty.If(!store.Wanted, vecty.Text("Members have been decided!")),
			),
			elem.UnorderedList(
				vecty.Markup(
					vecty.Class("list-group"),
				),
				items,
			),
		),
	)
}
