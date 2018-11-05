package route

import (
	"PKGAPPNAME/controllers/sample"

	"gopkg.in/sungora/app.v1/core"
)

func init() {
	core.SetRoute("/", &sample.ControlSample{})

	// sample group route
	core.SetRouteGroup2(map[string]map[string]core.ControllerFace{
		"/api/v1": map[string]core.ControllerFace{
			"/page1": &sample.ControlSample{},
			"/page2": &sample.ControlSample{},
			"/page3": &sample.ControlSample{},
		},
		"/api/v2": map[string]core.ControllerFace{
			"/page1": &sample.ControlSample{},
			"/page2": &sample.ControlSample{},
			"/page3": &sample.ControlSample{},
		},
	})
}
