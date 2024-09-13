package api

import "net/http"

func (self Router) CommonRoutes() {
	self.Router.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			println(err.Error())
		}
	})
}
