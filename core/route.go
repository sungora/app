package core

import (
	"errors"
)

type routesTyp map[string]ControllerFace

var Route = make(routesTyp)

func (r routesTyp) SetRoute(uri string, control ControllerFace) {
	r[uri] = control
}

func (r routesTyp) SetRouteGroup1(route map[string]ControllerFace) {
	for u, c := range route {
		r[u] = c
	}
}

func (r routesTyp) SetRouteGroup2(route map[string]map[string]ControllerFace) {
	for u1, map1 := range route {
		for u2, c := range map1 {
			r[u1+u2] = c
		}
	}
}

func (r routesTyp) SetRouteGroup3(route map[string]map[string]map[string]ControllerFace) {
	for u1, map1 := range route {
		for u2, map2 := range map1 {
			for u3, c := range map2 {
				r[u1+u2+u3] = c
			}
		}
	}
}

func (r routesTyp) GetRoute(uri string) (control ControllerFace, err error) {
	if _, ok := r[uri]; ok {
		return r[uri], nil
	}
	return nil, errors.New("controller not found from uri: " + uri)

}
