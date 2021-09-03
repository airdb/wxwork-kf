package middleware

import (
	"net/http"

	"github.com/airdb/sailor/version"
)

func myMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(version.GetBuildInfo().ToString()))
		// Here we are pssing our custom response writer to the next http handler.
		next.ServeHTTP(w, r)

		// Here we are adding our custom stuff to the response, which we received after http handler execution.
		// myResponseWriter.buf.WriteString(" and some additional modifications")
	})
}
