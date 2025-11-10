package handler

import (
	"net/http"

	"backend-fileprocessing/serverless"
)

// Handler é o entrypoint requisitado pela Vercel.
func Handler(w http.ResponseWriter, r *http.Request) {
	serverless.NewHandler().ServeHTTP(w, r)
}
