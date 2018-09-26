package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"honnef.co/go/js/dom"

	"github.com/nobonobo/gophertank/app/components"
	_ "github.com/nobonobo/gophertank/app/mutations"
	"github.com/nobonobo/gophertank/app/schema"
	_ "github.com/nobonobo/gophertank/app/services"
	"github.com/nobonobo/gophertank/app/store"
)

var (
	document     = dom.WrapDocument(js.Global.Get("document"))
	localStorage = js.Global.Get("localStorage")
)

const DEBUG = false

func main() {
	if localStorage.Get("uuid") == js.Undefined {
		localStorage.Set("uuid", uuid.New().String())
	}
	id := localStorage.Get("uuid").String()
	log.Print("uuid:", id)
	if localStorage.Get("name") == js.Undefined {
		localStorage.Set("name", "No Name")
	}
	name := localStorage.Get("name").String()
	log.Print("name:", name)
	document.AddEventListener("DOMContentLoaded", false, func(dom.Event) {
		if DEBUG {
			store.CurrentRoomMembers = []schema.Identity{
				{UUID: "", Name: "NoboNobo"},
				{UUID: "", Name: "さんぷる"},
			}
			v := &components.PlayView{}
			vecty.RenderBody(v)
		} else {
			vecty.RenderBody(&components.EntranceView{Name: name})
		}
	})
}
