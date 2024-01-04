package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/raulsilva-tech/Address-Finder/internal/infra/webserver/handlers"
)


func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	ah := handlers.NewAddressHandler()

	r.Get("/address/{cep}", ah.GetFastestAddressAnswer)

	http.ListenAndServe(":8888", r)
}
