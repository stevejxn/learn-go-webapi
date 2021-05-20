package handlers

import (
	"net/http"
)

func API() http.Handler {
	api := http.NewServeMux()
	api.HandleFunc("/status/ping", statusPingHandler)
	return api
}
