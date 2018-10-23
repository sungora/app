package core

import (
	"errors"
)

var routes = make(map[string]ControllerFace)

func SetRoute(uri string, control ControllerFace) {
	routes[uri] = control
}

func SetRouteGroup1(route map[string]ControllerFace) {
	for u, c := range route {
		routes[u] = c
	}
}

func SetRouteGroup2(route map[string]map[string]ControllerFace) {
	for u1, map1 := range route {
		for u2, c := range map1 {
			routes[u1+u2] = c
		}
	}
}

func SetRouteGroup3(route map[string]map[string]map[string]ControllerFace) {
	for u1, map1 := range route {
		for u2, map2 := range map1 {
			for u3, c := range map2 {
				routes[u1+u2+u3] = c
			}
		}
	}
}

func GetRoute(uri string) (control ControllerFace, err error) {
	if _, ok := routes[uri]; ok {
		return routes[uri], nil
	}
	return nil, errors.New("route not found: " + uri)

}
