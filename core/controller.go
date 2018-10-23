package core

import (
	"net/http"
)

// ContraFace is an interface to uniform all controller handler.
type ControllerFace interface {
	Init(w http.ResponseWriter, r *http.Request)
	GET() (err error)
	POST() (err error)
	PUT() (err error)
	DELETE() (err error)
	OPTIONS() (err error)
	Render()
}

