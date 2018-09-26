package mutations

import (
	"html"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"reflect"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/websocket"
	"honnef.co/go/js/dom"

	"github.com/nobonobo/gophertank/app/actions"
	"github.com/nobonobo/gophertank/app/components"
	"github.com/nobonobo/gophertank/app/dispatcher"
	"github.com/nobonobo/gophertank/app/schema"
	"github.com/nobonobo/gophertank/app/services"
	"github.com/nobonobo/gophertank/app/store"
)

var (
	localStorage = js.Global.Get("localStorage")
	location     = js.Global.Get("location")
	conn         net.Conn
	lastView     vecty.Component
)

func init() {
	dispatcher.Register(onAction)
}

func onAction(action schema.Action) {
	switch a := action.(type) {
	default:
		log.Println("unknown action:", action)
	case *actions.Enter:
		log.Print("enter:", a.Name)
		localStorage.Set("name", a.Name)
		store.Wanted = true
		store.Identity = schema.Identity{
			Name: html.EscapeString(a.Name),
			UUID: localStorage.Get("uuid").String(),
		}
		lastView = &components.RoomView{}
		vecty.RenderBody(lastView)
		go func() {
			uri := "ws://"
			if location.Get("protocol").String() == "https:" {
				uri = "wss://"
			}
			uri += location.Get("host").String() + "/ws/"
			ws, err := websocket.Dial(uri)
			if err != nil {
				dom.GetWindow().Alert("Connection Failed: " + err.Error())
				dispatcher.Dispatch(&actions.Leave{})
				return
			}
			conn = ws
			log.Print("websocket connected")
			rpc.ServeCodec(jsonrpc.NewServerCodec(ws))
			log.Print("websocket disconnected")
		}()
	case *actions.Leave:
		log.Print("leave:", store.Identity.Name)
		if dom.GetWindow().Confirm("Leave from this room?") {
			if conn != nil {
				conn.Close()
				conn = nil
			}
			vecty.RenderBody(&components.EntranceView{Name: store.Identity.Name})
			store.CurrentRoomMembers = store.CurrentRoomMembers[0:0]
		}
	case *actions.UpdateMembers:
		log.Print("update member:", len(a.Members))
		if !reflect.DeepEqual(store.CurrentRoomMembers, a.Members) {
			store.CurrentRoomMembers = a.Members
			store.Wanted = len(a.Members) < store.MaxMembers
			vecty.Rerender(lastView)
		}
	case *actions.Begin:
		log.Print("preperation begin member:", len(a.Members))
		store.Wanted = false
		store.CurrentRoomMembers = a.Members
		vecty.Rerender(lastView)
	case *actions.Abort:
		log.Print("abort")
		dom.GetWindow().Alert("P2P connection failed")
		vecty.RenderBody(&components.EntranceView{Name: store.Identity.Name})
		store.CurrentRoomMembers = store.CurrentRoomMembers[0:0]
	case *actions.CreateOffer:
		log.Println("create offer")
		peerConnection := js.Global.Get("RTCPeerConnection").New(js.M{
			"iceServers": js.S{js.M{"urls": "stun:stun.l.google.com:19302"}},
		})
		peerConnection.Set("onicecandidate", func(ev *js.Object) {
			if ev.Get("candidate") == nil {
				a.Response <- peerConnection.Get("localDescription").Get("sdp").String()
			}
		})
		dc := peerConnection.Call("createDataChannel", "rpc")
		peerConnection.Call("createOffer").Call(
			"then", func(sdp *js.Object) {
				peerConnection.Call("setLocalDescription", sdp)
			},
		)
		store.Others[a.UUID] = &store.Client{
			Client:         jsonrpc.NewClient(services.NewDCConn(dc)),
			PeerConnection: peerConnection,
		}
	case *actions.CreateAnswer:
		log.Println("create answer")
		peerConnection := js.Global.Get("RTCPeerConnection").New(js.M{
			"iceServers": js.S{js.M{"urls": "stun:stun.l.google.com:19302"}},
		})
		peerConnection.Set("ondatachannel", func(ev *js.Object) {
			c := services.NewDCConn(ev.Get("channel"))
			c.Set("onopen", func() {
				go func() {
					rpc.ServeCodec(jsonrpc.NewServerCodec(c))
				}()
			})
		})
		sdp := js.Global.Get("RTCSessionDescription").New(js.M{
			"type": "offer",
			"sdp":  a.SDP,
		})
		peerConnection.Set("onicecandidate", func(ev *js.Object) {
			if ev.Get("candidate") == nil {
				a.Response <- peerConnection.Get("localDescription").Get("sdp").String()
			}
		})
		peerConnection.Call("setRemoteDescription", sdp).Call(
			"then", func() {
				peerConnection.Call("createAnswer").Call(
					"then", func(sdp *js.Object) {
						peerConnection.Call("setLocalDescription", sdp)
					},
				)
			},
		)
	case *actions.CreateConn:
		log.Println("create connection")
		client := store.Others[a.UUID]
		if client != nil {
			sdp := js.Global.Get("RTCSessionDescription").New(js.M{
				"type": "answer",
				"sdp":  a.SDP,
			})
			client.PeerConnection.Call("setRemoteDescription", sdp)
		}
	case *actions.End:
		log.Println("preperation end")
		lastView = &components.PlayView{}
		vecty.RenderBody(lastView)
		go func() {
			for id, c := range store.Others {
				var identity schema.Identity
				if err := c.Call("Player.GetIdentity", schema.None, &identity); err != nil {
					log.Println(err)
				}
				log.Println(identity.Name, id == identity.UUID)
			}
		}()
	}
}
