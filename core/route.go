package core

import (
	"errors"
)

type routesTyp map[string]ControllerFace

var Route = make(routesTyp)

func (self routesTyp) SetRoute(uri string, control ControllerFace) {
	self[uri] = control
}

func (self routesTyp) SetRouteGroup1(route map[string]ControllerFace) {
	for u, c := range route {
		self[u] = c
	}
}

func (self routesTyp) SetRouteGroup2(route map[string]map[string]ControllerFace) {
	for u1, map1 := range route {
		for u2, c := range map1 {
			self[u1+u2] = c
		}
	}
}

func (self routesTyp) SetRouteGroup3(route map[string]map[string]map[string]ControllerFace) {
	for u1, map1 := range route {
		for u2, map2 := range map1 {
			for u3, c := range map2 {
				self[u1+u2+u3] = c
			}
		}
	}
}

func (self routesTyp) GetRoute(uri string) (control ControllerFace, err error) {
	if _, ok := self[uri]; ok {
		return self[uri], nil
	}
	return nil, errors.New("controller not found from uri: " + uri)

}
