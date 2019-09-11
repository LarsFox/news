package api

import (
	"net/http"
)

// route описывает поля обработчика запроса. Используется при добавлении путей
// и обработчиков запросов в мультиплексов (mux).
type route struct {
	Method   string
	Path     string
	Name     string
	HndlrGen func() http.Handler
}

// addHandlers добавляет пути и обработчики запросов в мультиплексор (mux).
func (m *Manager) addHandlers(routes []route) {
	for _, r := range routes {
		m.router.Methods(r.Method).Path(r.Path).Name(r.Name).Handler(r.HndlrGen())
	}
}
