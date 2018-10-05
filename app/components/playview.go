package components

import (
	"log"
	"math"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	"github.com/lngramos/three"
	"github.com/vecty/vthree"

	"github.com/nobonobo/gophertank/app/store"
)

// PlayView ...
type PlayView struct {
	vecty.Core
	scene     *three.Scene
	camera    three.PerspectiveCamera
	renderer  *three.WebGLRenderer
	mesh      *three.Mesh
	stick     *js.Object
	mouse     *js.Object
	device    *js.Object
	person    *js.Object
	velocityX float64
	velocityY float64
	lastTime  float64
}

// Render ...
func (p *PlayView) Render() vecty.ComponentOrHTML {
	vecty.SetTitle("Play")
	return elem.Body(
		elem.Div(
			vecty.Markup(
				prop.ID("stick-area"),
				vecty.Class("w-25", "h-100", "position-absolute"),
			),
		),
		vthree.WebGLRenderer(vthree.WebGLOptions{
			Init:     p.init,
			Shutdown: p.shutdown,
		}),
		elem.Canvas(
			vecty.Markup(
				prop.ID("ids"),
				vecty.Style("width", "320px"),
				vecty.Style("height", "160px"),
				vecty.Style("display", "none"),
			),
		),
	)
}

func (p *PlayView) shutdown(renderer *three.WebGLRenderer) {
	// After shutdown, we shouldn't use any of these anymore.
	p.scene = nil
	p.camera = three.PerspectiveCamera{}
	p.person = nil
	p.renderer = nil
	p.mesh = nil
	p.stick = nil
	js.Global.Get("bodyScrollLock").Call("clearAllBodyScrollLocks")
}

func (p *PlayView) onResize() {
	time.AfterFunc(500*time.Millisecond, func() {
		windowWidth := js.Global.Get("innerWidth").Float()
		windowHeight := js.Global.Get("innerHeight").Float()
		devicePixelRatio := js.Global.Get("devicePixelRatio").Float()
		p.camera.Set("aspect", windowWidth/windowHeight)
		p.renderer.SetPixelRatio(devicePixelRatio)
		p.renderer.SetSize(windowWidth, windowHeight, true)
		p.camera.Call("updateProjectionMatrix")
	})
}

func (p *PlayView) init(renderer *three.WebGLRenderer) {
	js.Global.Get("bodyScrollLock").Call("disableBodyScroll", renderer.Object)
	p.stick = js.Global.Get("VirtualJoystick").New(js.M{
		"container":        js.Global.Get("document").Call("getElementById", "stick-area"),
		"limitStickTravel": true,
		"stickRadius":      50,
		"mouseSupport":     true,
	})

	canvas := js.Global.Get("document").Call("getElementById", "ids").Call("getContext", "2d")
	canvas.Set("font", "18px kosugi maru bold")
	canvas.Set("textAlign", "center")
	canvas.Set("fillStyle", "#ffffffe0")
	canvas.Set("shadowColor", "black")
	canvas.Set("shadowOffsetX ", 0)
	canvas.Set("shadowOffsetY ", 0)
	canvas.Set("shadowBlur  ", 8)
	for i, m := range store.CurrentRoomMembers {
		canvas.Call("strokeText", m.Name, 160+1, (i+1)*20+1, 320)
		canvas.Call("strokeText", m.Name, 160-1, (i+1)*20+1, 320)
		canvas.Call("strokeText", m.Name, 160+1, (i+1)*20-1, 320)
		canvas.Call("strokeText", m.Name, 160-1, (i+1)*20-1, 320)
		canvas.Call("fillText", m.Name, 160, (i+1)*20, 320)
	}

	p.renderer = renderer
	windowWidth := js.Global.Get("innerWidth").Float()
	windowHeight := js.Global.Get("innerHeight").Float()
	devicePixelRatio := js.Global.Get("devicePixelRatio").Float()

	p.camera = three.NewPerspectiveCamera(75, windowWidth/windowHeight, 0.5, 1000)
	p.camera.Position.Set(0.0, 0.5, 1.0)

	js.Global.Call("addEventListener", "resize", p.onResize)
	p.device = js.Global.Get("THREE").Get("DeviceOrientationControls").New(p.camera.Object)
	p.camera.Call("lookAt", three.NewVector3(0, 0, 0))
	p.scene = three.NewScene()
	p.scene.Set("fog", js.Global.Get("THREE").Get("FogExp2").New(0xffffff, 0.015))

	pcMode := true

	js.Global.Call("addEventListener", "deviceorientation", func(ev *js.Object) {
		if ev.Get("alpha") != nil {
			pcMode = false
		}
	})
	time.AfterFunc(0*time.Millisecond, func() {
		if !pcMode {
			log.Println("smart phone mode")
			p.person = p.camera.Object
			p.person.Call("translateZ", 3)
			//p.scene.Call("add", p.person)
		} else {
			log.Println("pc mouse mode")
			p.device = nil
			// input for PC
			p.mouse = js.Global.Get("THREE").Get("PointerLockControls").New(p.camera.Object)
			renderer.Get("domElement").Call("addEventListener", "click", func(ev *js.Object) {
				p.mouse.Call("lock")
			})
			js.Global.Get("document").Call("addEventListener", "keydown", func(ev *js.Object) {
				switch ev.Get("keyCode").Int() {
				case 38, 87: // up or w
					p.velocityY = -1.0
				case 40, 83: // down or s
					p.velocityY = 1.0
				case 37, 65: // left or a
					p.velocityX = -1.0
				case 39, 68: // right or d
					p.velocityX = +1.0
				}
			})
			js.Global.Get("document").Call("addEventListener", "keyup", func(ev *js.Object) {
				switch ev.Get("keyCode").Int() {
				case 38, 87: // up or w
					if p.velocityY < 0 {
						p.velocityY = 0.0
					}
				case 40, 83: // down or s
					if p.velocityY > 0 {
						p.velocityY = 0.0
					}
				case 37, 65: // left or a
					if p.velocityX < 0 {
						p.velocityX = 0.0
					}
				case 39, 68: // right or d
					if p.velocityX > 0 {
						p.velocityX = 0.0
					}
				}
			})
			p.person = p.mouse.Call("getObject")
			p.person.Call("translateZ", 5)
			p.scene.Call("add", p.person)
		}
		// Begin animating.
		p.animate()
	})

	light := three.NewDirectionalLight(three.NewColor(255, 255, 255), 0.8)
	//light.Position.Set(256, 256, 256).Normalize()
	p.scene.Add(light)

	p.renderer.SetPixelRatio(devicePixelRatio)
	p.renderer.SetSize(windowWidth, windowHeight, true)

	loader := js.Global.Get("THREE").Get("TextureLoader").New()
	waterGeometry := js.Global.Get("THREE").Get("PlaneBufferGeometry").New(10000, 10000)
	water := js.Global.Get("THREE").Get("Water").New(
		waterGeometry,
		js.M{
			"textureWidth":  512,
			"textureHeight": 512,
			"waterNormals": loader.Call(
				"load", "textures/waternormals.jpg",
				func(texture *js.Object) {
					texture.Set("wrapS", js.Global.Get("THREE").Get("RepeatWrapping"))
					texture.Set("wrapT", js.Global.Get("THREE").Get("RepeatWrapping"))
				},
			),
			"alpha":           1.0,
			"sunDirection":    light.Get("position").Call("clone").Call("normalize"),
			"sunColor":        0xffffff,
			"waterColor":      0x001e0f,
			"distortionScale": 3.7,
			"fog":             p.scene.Get("fog") != js.Undefined,
		},
	)
	water.Get("rotation").Set("x", -math.Pi/2)
	p.scene.Call("add", water)

	//===========
	sky := js.Global.Get("THREE").Get("Sky").New()
	sky.Get("scale").Call("setScalar", 10000)
	p.scene.Call("add", sky)
	uniforms := sky.Get("material").Get("uniforms")
	uniforms.Get("turbidity").Set("value", 10)
	uniforms.Get("rayleigh").Set("value", 2)
	uniforms.Get("luminance").Set("value", 1)
	uniforms.Get("mieCoefficient").Set("value", 0.005)
	uniforms.Get("mieDirectionalG").Set("value", 0.8)
	var parameters = js.M{
		"distance":    400.0,
		"inclination": 0.49,
		"azimuth":     0.205,
	}
	cubeCamera := js.Global.Get("THREE").Get("CubeCamera").New(1, 20000, 256)
	cubeCamera.Get("renderTarget").Get("texture").Set(
		"minFilter", js.Global.Get("THREE").Get("LinearMipMapLinearFilter"))
	updateSun := func() {
		theta := math.Pi * (parameters["inclination"].(float64) - 0.5)
		phi := 2 * math.Pi * (parameters["azimuth"].(float64) - 0.5)
		light.Get("position").Set("x", parameters["distance"].(float64)*math.Cos(phi))
		light.Get("position").Set("y", parameters["distance"].(float64)*math.Sin(phi)*math.Sin(theta))
		light.Get("position").Set("z", parameters["distance"].(float64)*math.Sin(phi)*math.Cos(theta))
		sky.Get("material").Get("uniforms").Get("sunPosition").Set("value", light.Get("position").Call("copy", light.Get("position")))
		water.Get("material").Get("uniforms").Get("sunDirection").Get("value").Call("copy", light.Get("position")).Call("normalize")
		cubeCamera.Call("update", renderer.Object, p.scene.Object)
	}
	updateSun()

	// Create cube
	geometry := three.NewBoxGeometry(&three.BoxGeometryParameters{
		Width:  1,
		Height: 1,
		Depth:  1,
	})

	// geometry2 := three.NewCircleGeometry(three.CircleGeometryParameters{
	// 	Radius:      50,
	// 	Segments:    20,
	// 	ThetaStart:  0,
	// 	ThetaLength: 2,
	// })

	materialParams := three.NewMaterialParameters()
	materialParams.Color = three.NewColor(0, 123, 211)
	materialParams.Shading = three.SmoothShading
	materialParams.Side = three.FrontSide
	material := three.NewMeshBasicMaterial(materialParams)
	// material := three.NewMeshLambertMaterial(materialParams)
	// material := three.NewMeshPhongMaterial(materialParams)
	p.mesh = three.NewMesh(geometry, material)
	p.scene.Add(p.mesh)

	p.lastTime = js.Global.Get("performance").Call("now").Float()
}

func (p *PlayView) animate() {
	if p.renderer == nil {
		// We shutdown, stop animation.
		return
	}
	tm := js.Global.Get("performance").Call("now").Float()
	defer func() { p.lastTime = tm }()
	delta := tm - p.lastTime
	js.Global.Call("requestAnimationFrame", p.animate)
	p.renderer.Render(p.scene, p.camera)
	if p.device != nil {
		p.device.Call("update")
	}
	v := p.stick.Call("deltaX").Float()/50.0 + p.velocityX
	w := p.stick.Call("deltaY").Float()/50.0 + p.velocityY
	if p.person != nil {
		p.person.Call("translateX", 0.01*v*delta)
		p.person.Call("translateZ", 0.01*w*delta)
	}
}
