package handlers

import (
	"fmt"
	"github.com/dimfeld/httptreemux/v5"
	"net/http"
	"time"
)

func statusPingHandler(w http.ResponseWriter, r *http.Request) {
	ctxData := httptreemux.ContextData(r.Context())
	routePath := ctxData.Route()

	pingResponse := fmt.Sprintf("uri: %s - pong: %v",
		routePath,
		time.Now().UTC().Format(time.RFC3339))

	w.Write([]byte(pingResponse))
}
