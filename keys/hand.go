package keys

func init() {
	Handler = &handler{
		Session:      "Session",
		Status:       "Status",
		RoutePattern: "RoutePattern",
		Log:          "LogHandler",
	}
}

var Handler *handler

type handler struct {
	Session      string
	Status       string
	RoutePattern string
	Log          string
}
