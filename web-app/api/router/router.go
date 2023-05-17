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

	r.Get("/livez", health.Read)

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		http.HandlerFunc()
		r.Method("GET", "/balance", requestlog.NewHandler(bookAPI.List, l))
		r.Method("POST", "/transfer", requestlog.NewHandler(bookAPI.Create, l))
		r.Method("POST", "/register", requestlog.NewHandler(bookAPI.Read, l))
		r.Method("POST", "/authentication", requestlog.NewHandler(bookAPI.Update, l))

		// use container DI like uber-go/dig instead of manual
		DIBalanceHandler := func(w http.ResponseWriter, r *http.Request) {
			handlers.BalanceHandler(w, r, connector)
		}
		DIAuthHandler := func(w http.ResponseWriter, r *http.Request) { handlers.AuthHandler(w, r, connector) }
		DIRegHandler := func(w http.ResponseWriter, r *http.Request) { handlers.RegHandler(w, r, connector) }
		DITransferHandler := func(w http.ResponseWriter, r *http.Request) {
			handlers.TransferHandler(w, r, kafkaW)
		}

		http.HandleFunc("/balance", middleware.AuthMiddleware(DIBalanceHandler))
		http.HandleFunc("/authentication", DIAuthHandler)
		http.HandleFunc("/register", DIRegHandler)
		http.HandleFunc("/transfer", middleware.AuthMiddleware(DITransferHandler))
	})

	return r
}
