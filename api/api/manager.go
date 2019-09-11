package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/LarsFox/news/api/queues"
)

// Manager ...
type Manager struct {
	router  *mux.Router
	queuesM *queues.Manager
}

func (m *Manager) Listen(addr string) error {
	log.Println("Server started at", addr)
	srv := &http.Server{
		Addr:    addr,
		Handler: m.router,
	}
	return srv.ListenAndServe()
}

func NewManager(queuesM *queues.Manager) *Manager {
	m := &Manager{
		router:  mux.NewRouter(),
		queuesM: queuesM,
	}

	m.addHandlers([]route{
		{ // Получить новость.
			Method:   "GET",
			Path:     "/news/{news_id}/",
			Name:     "GetNewsPiece",
			HndlrGen: m.hGetNewsPiece,
		},
	})

	return m
}

// send отправляет клиенту ответ.
func (m *Manager) send(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	resp := map[string]interface{}{
		"ok":     true,
		"result": data,
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		log.Println(err)
		return
	}
}

// sendErr отправляет клиенту сообщение об ошибке.
func (m *Manager) sendErr(w http.ResponseWriter, r *http.Request, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	resp := map[string]interface{}{
		"ok":         false,
		"error_code": code,
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		log.Println(err)
		return
	}
}

func (m *Manager) hGetNewsPiece() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		newsID := vars["news_id"]
		reply, err := m.queuesM.RequestNewsPiece(newsID)
		if err != nil {
			m.sendErr(w, r, 500)
			return
		}
		switch reply.ErrorCode {
		case 0:
		case 1:
			m.sendErr(w, r, 500)
			return
		case 1001:
			m.sendErr(w, r, 404)
			return
		default:
			log.Printf("unknown error code: %d", reply.ErrorCode)
			m.sendErr(w, r, 500)
			return
		}

		m.send(w, r, map[string]interface{}{
			"header": reply.Header,
			"date":   reply.Date,
		})
	})
}
