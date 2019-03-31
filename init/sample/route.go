package sample

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	middlewareChi "github.com/go-chi/chi/middleware"

	"github.com/sungora/app/internal/core"
	"github.com/sungora/app/internal/sample/apiv1/groups"
	"github.com/sungora/app/internal/sample/apiv1/users"
	"github.com/sungora/app/internal/sample/middleware"
	"github.com/sungora/app/internal/sample/page"
	"github.com/sungora/app/servhttp/middlew"
)

func routes(r *chi.Mux) {

	r.NotFound(middlew.NotFound)
	r.Use(middlew.TimeoutContext(time.Second * time.Duration(core.Cfg.Http.WriteTimeout-1)))
	r.Use(middlewareChi.Recoverer)
	r.Use(middlewareChi.Logger)
	r.Use(middleware.SampleRoot)

	// Group 1
	r.Group(func(r chi.Router) {
		r.HandleFunc("/", page.Main)
		r.HandleFunc("/api", page.Api)
		r.HandleFunc("/api/v1", page.ApiV1)
	})

	// Group 2
	r.Group(func(r chi.Router) {
		r.Use(middleware.SamplePing)
		r.Get("/ping", page.Ping)                                          // sample more routes
		r.Get("/version", page.Version)                                    // sample more routes
		r.Get("/test/{testID}/order/{orderID}/page/{pageID}", page.Sample) // sample more routes
	})

	r.Mount("/api/v1/sample", sampleApiV1())
}

// sampleApiV1 роутинг api запросов /api/v1/sample
func sampleApiV1() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.SampleOne)
	r.HandleFunc("/", page.ApiV1)
	// sample routes for "users" resource
	r.Route("/user", func(r chi.Router) {
		r.Use(middleware.SampleTwo)
		r.Get("/", users.Gets)        // array /user
		r.Get("/search", users.Gets)  // array search /user
		r.Post("/", users.Post)       // POST /user
		r.Options("/", users.Options) // OPTIONS /user
		r.Route("/{ID}", func(r chi.Router) {
			r.Get("/", users.Get)       // GET /user/123
			r.Put("/", users.Put)       // PUT /user/123
			r.Delete("/", users.Delete) // DELETE /user/123
		})

	})
	// sample routes for "groups" resource
	r.Route("/group", func(r chi.Router) {
		r.Use(middleware.SampleFour)
		r.Get("/", groups.Gets)        // array /group
		r.Get("/search", groups.Gets)  // array search /group
		r.Post("/", groups.Post)       // POST /group
		r.Options("/", groups.Options) // OPTIONS /group
		r.Route("/{ID}", func(r chi.Router) {
			r.Get("/", groups.Get)       // GET /group/123
			r.Put("/", groups.Put)       // PUT /group/123
			r.Delete("/", groups.Delete) // DELETE /group/123
		})

	})
	return r
}
