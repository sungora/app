package app

import (
	"errors"
)

type routesTyp map[string]func() ControllerFace

var Route = make(routesTyp)

func (r routesTyp) Set(path string, f func() ControllerFace) {
	r[path] = f
}

func (r routesTyp) Get(path string) (control ControllerFace, err error) {
	if _, ok := r[path]; ok {
		return r[path](), nil
	}
	return nil, errors.New("controller not found from uri: " + path)

}

func (r routesTyp) Path(pathSegment string) routePath {
	return routePath(pathSegment)
}

type routePath string

func (r routePath) Path(pathSegment string) routePath {
	return r + routePath(pathSegment)
}

func (r routePath) Set(pathSegment string, f func() ControllerFace) routePath {
	Route[string(r)+pathSegment] = f
	return r
}
