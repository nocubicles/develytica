package middleware

import "net/http"

type AppHandler func(http.ResponseWriter, *http.Request) (int, error)

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := fn(w, r); err != nil {
		switch status {
		case http.StatusNotFound:
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
