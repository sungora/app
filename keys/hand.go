package keys

func init() {
	Hand = &hand{
		Session:      "Session",
		Status:       "Status",
		RoutePattern: "RoutePattern",
	}
}

var Hand *hand

type hand struct {
	Session      string
	Status       string
	RoutePattern string
}
