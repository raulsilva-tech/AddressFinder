package main

import (
	"fmt"
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
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It's running"))
	})
	fmt.Println("Serving on port 9090")
	err := http.ListenAndServe(":9090", r)
	if err != nil {
		panic(err)
	}
}
