package router

import (
	"net/http"

	"github.com/JakubDaleki/transfer-app/webapp/api/handlers"
	"github.com/JakubDaleki/transfer-app/webapp/api/router/middleware"
	"github.com/JakubDaleki/transfer-app/webapp/utils/db"
	"github.com/go-chi/chi/v5"
	"github.com/segmentio/kafka-go"
)

func New(connector *db.Connector, kafkaW *kafka.Writer) *chi.Mux {
	r := chi.NewRouter()

	DIAuthHandler := func(w http.ResponseWriter, r *http.Request) { handlers.AuthHandler(w, r, connector) }
	DIRegHandler := func(w http.ResponseWriter, r *http.Request) { handlers.RegHandler(w, r, connector) }
	r.Post("/register", http.HandlerFunc(DIRegHandler))
	r.Post("/authentication", http.HandlerFunc(DIAuthHandler))

	r.Route("/account", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// use container DI like uber-go/dig instead of manual
		DIBalanceHandler := func(w http.ResponseWriter, r *http.Request) { handlers.BalanceHandler(w, r, connector) }
		DITransferHandler := func(w http.ResponseWriter, r *http.Request) { handlers.TransferHandler(w, r, kafkaW) }
		r.Method("GET", "/balance", http.HandlerFunc(DIBalanceHandler))
		r.Method("POST", "/transfer", http.HandlerFunc(DITransferHandler))

	})

	return r
}
