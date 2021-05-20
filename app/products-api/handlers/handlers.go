package handlers

import (
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
)

func API() http.Handler {
	router := httptreemux.NewContextMux()
	group := router.NewGroup("/api")
	group.GET("/status/ping", statusPingHandler)
	return router
}
